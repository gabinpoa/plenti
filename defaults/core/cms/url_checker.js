export const makeUrl = string => {
  try {
    return new URL(string);
  } catch (_) {
    return new URL("https://gitlab.com");
  }
}