/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./presentation/**/*.templ", "./components/**/*.templ"],
  theme: {
    extend: {},
    container: {
      screens: {
        sm: "600px",
        md: "728px",
        lg: "1020px",
        xl: "1280px",
        "2xl": "1280px",
      },
    },
  },
  plugins: [],
};
