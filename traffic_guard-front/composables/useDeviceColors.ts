export const useDeviceColors = () => {
  const generateColor = (deviceName: string) => {
    let hash = 0;
    for (let i = 0; i < deviceName.length; i++) {
      hash = deviceName.charCodeAt(i) + ((hash << 5) - hash);
    }
    
    let hue = Math.abs(hash) % 360;
    
    if (hue >= 300) {
      hue = (hue - 300) * (60 / 60);
    }
    
    const saturation = 70 + (hash % 20);
    const lightness = 85 + (hash % 10);
    const textLightness = 35 + (hash % 15);
    
    return {
      background: `hsl(${hue}, ${saturation}%, ${lightness}%)`,
      color: `hsl(${hue}, ${saturation}%, ${textLightness}%)`
    };
  };

  return { generateColor };
};