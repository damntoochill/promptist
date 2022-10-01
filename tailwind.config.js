/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./templates/*.tmpl"],
  theme: {
    fontFamily: {
      sans: ["-apple-system", "BlinkMacSystemFont", "Segoe UI",
      "Roboto", "Oxygen-Sans", "Ubuntu", "Cantarell",
      "Helvetica Neue", "sans-serif"],
      serif: ['Iowan Old Style', 'Apple Garamond', 'Baskerville', 'Times New Roman', 'Droid Serif', 'Times', 'Source Serif Pro', 'serif', 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol'],
      mono: ['SFMono-Regular', 'Menlo', 'Monaco','Consolas', "Liberation Mono", "Courier New", 'monospace'],
    },
    extend: {},
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
  ]
}
