export const useRootFontSize = () => {
  const rootFontSize = ref(0);

  const update = () => {
    rootFontSize.value = parseFloat(
      getComputedStyle(document.documentElement).fontSize
    );
  };

  onMounted(() => {
    update();
    window.addEventListener('resize', update);
  });

  onBeforeUnmount(() => {
    window.removeEventListener('resize', update);
  });

  return rootFontSize;
};
