/**
 * Encode file path for URL query parameter
 */
export function encodePath(path: string): string {
  return encodeURIComponent(path);
}

/**
 * Decode file path from URL query parameter
 */
export function decodePath(encodedPath: string): string {
  try {
    return decodeURIComponent(encodedPath);
  } catch (e) {
    console.error('Failed to decode path:', e);
    return encodedPath;
  }
}

