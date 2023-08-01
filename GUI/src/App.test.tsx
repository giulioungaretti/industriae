import { test } from "@jest/globals"
import { render, screen } from "@testing-library/react"
import App from "./App"

test("renders learn react link", () => {
  // RENDER APP
  render(<App />)

  expect(screen.getByText(/Pressure/i)).toBeInTheDocument()
})
