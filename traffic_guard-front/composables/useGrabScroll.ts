import type { Ref } from 'vue';

const DRAG_THRESHOLD_PX = 5;

const POINTER_MOVE_OPTS: AddEventListenerOptions = { passive: false };

function isGrabIgnoredTarget(target: EventTarget | null): boolean {
  if (!target || !(target instanceof Element)) return false;
  return Boolean(
    target.closest(
      [
        'button',
        'a',
        'input',
        'textarea',
        'select',
        'option',
        '[role="button"]',
        '[contenteditable="true"]',
        'canvas',
        'svg',
        'iframe',
        'label',
        'summary',
        '[data-grab-scroll-ignore]',
      ].join(', '),
    ),
  );
}

export function useGrabScroll(scrollRef: Ref<HTMLElement | null>) {
  let pointerDown = false;
  let dragStarted = false;
  let startX = 0;
  let startY = 0;
  let startScrollLeft = 0;
  let startScrollTop = 0;
  let rafId: number | undefined;

  function updateGrabCursor(el: HTMLElement) {
    const x = el.scrollWidth > el.clientWidth + 1;
    const y = el.scrollHeight > el.clientHeight + 1;
    el.classList.toggle('grab-scroll--can-grab', x || y);
  }

  let ro: ResizeObserver | undefined;
  function observeScrollAndContent(el: HTMLElement) {
    ro = new ResizeObserver(() => updateGrabCursor(el));
    ro.observe(el);
    const child = el.firstElementChild;
    if (child instanceof HTMLElement) {
      ro.observe(child);
    }
  }

  function cancelPendingRaf() {
    if (rafId !== undefined) {
      cancelAnimationFrame(rafId);
      rafId = undefined;
    }
  }

  function cleanupDocumentListeners() {
    cancelPendingRaf();
    // document недоступен при SSR — watch(immediate) дёргает cleanup ещё на сервере
    if (typeof document === 'undefined') return;
    document.removeEventListener('pointermove', onPointerMove, POINTER_MOVE_OPTS);
    document.removeEventListener('pointerup', onPointerUp);
    document.removeEventListener('pointercancel', onPointerUp);
  }

  function onPointerMove(e: Event) {
    const ev = e as PointerEvent;
    const el = scrollRef.value;
    if (!el || !pointerDown) return;

    const dx = ev.clientX - startX;
    const dy = ev.clientY - startY;

    if (!dragStarted) {
      if (Math.hypot(dx, dy) < DRAG_THRESHOLD_PX) return;
      dragStarted = true;
      el.classList.add('grab-scroll--dragging');
    }

    ev.preventDefault();
    cancelPendingRaf();
    rafId = requestAnimationFrame(() => {
      el.scrollLeft = startScrollLeft - dx;
      el.scrollTop = startScrollTop - dy;
      rafId = undefined;
    });
  }

  function onPointerUp(e: PointerEvent) {
    cleanupDocumentListeners();

    const el = scrollRef.value;
    if (!pointerDown) return;
    pointerDown = false;

    if (dragStarted && el) {
      el.classList.remove('grab-scroll--dragging');
      try {
        el.releasePointerCapture(e.pointerId);
      } catch {
        /* ignore */
      }

      const suppressClick = (ev: Event) => {
        ev.preventDefault();
        ev.stopPropagation();
      };
      el.addEventListener('click', suppressClick, { capture: true, once: true });
    }

    dragStarted = false;
  }

  function onPointerDown(e: PointerEvent) {
    const el = scrollRef.value;
    if (!el || e.pointerType !== 'mouse' || e.button !== 0) return;
    if (isGrabIgnoredTarget(e.target)) return;

    pointerDown = true;
    dragStarted = false;
    startX = e.clientX;
    startY = e.clientY;
    startScrollLeft = el.scrollLeft;
    startScrollTop = el.scrollTop;

    try {
      el.setPointerCapture(e.pointerId);
    } catch {
      /* недоступен */
    }

    document.addEventListener('pointermove', onPointerMove, POINTER_MOVE_OPTS);
    document.addEventListener('pointerup', onPointerUp);
    document.addEventListener('pointercancel', onPointerUp);
  }

  function bind(el: HTMLElement | null) {
    ro?.disconnect();
    ro = undefined;

    if (!el) return;

    updateGrabCursor(el);
    observeScrollAndContent(el);

    const onScroll = () => updateGrabCursor(el);
    el.addEventListener('scroll', onScroll, { passive: true });

    el.addEventListener('pointerdown', onPointerDown);

    return () => {
      el.removeEventListener('pointerdown', onPointerDown);
      el.removeEventListener('scroll', onScroll);
      ro?.disconnect();
      ro = undefined;
    };
  }

  let unbind: (() => void) | undefined;

  watch(
    scrollRef,
    (el) => {
      unbind?.();
      unbind = undefined;
      pointerDown = false;
      dragStarted = false;
      cleanupDocumentListeners();
      unbind = bind(el) ?? undefined;
    },
    { immediate: true },
  );

  onUnmounted(() => {
    unbind?.();
    cleanupDocumentListeners();
  });
}
