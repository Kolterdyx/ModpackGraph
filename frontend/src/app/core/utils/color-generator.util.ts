/**
 * Generate a deterministic color from a string (mod ID)
 */
export function generateColor(str: string): string {
  let hash = 0;
  for (let i = 0; i < str.length; i++) {
    hash = str.charCodeAt(i) + ((hash << 5) - hash);
  }

  const hue = Math.abs(hash % 360);
  const saturation = 65 + (Math.abs(hash) % 20);
  const lightness = 50 + (Math.abs(hash >> 8) % 20);

  return `hsl(${hue}, ${saturation}%, ${lightness}%)`;
}

/**
 * Generate a color palette for a group of items
 */
export function generateColorPalette(items: string[]): Map<string, string> {
  const palette = new Map<string, string>();
  items.forEach((item) => {
    palette.set(item, generateColor(item));
  });
  return palette;
}

