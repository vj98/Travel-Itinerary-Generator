import React from "react";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import "@testing-library/jest-dom";
import TripPlanner from "@/app/dashboard/page";
import { callApi } from "@/app/utils/api-utils";
import { useRouter } from "next/navigation";

// Mock the dependencies
jest.mock("@/app/utils/api-utils");
jest.mock("next/router", () => ({
  useRouter: jest.fn(),
}));

const mockPush = jest.fn();

beforeEach(() => {
  useRouter.mockImplementation(() => ({
    push: mockPush,
    prefetch: jest.fn().mockResolvedValue(undefined),
    route: "/",
    pathname: "/",
    query: {},
    asPath: "/",
  }));
});

describe("TripPlanner", () => {
  test("renders TripPlanner component", () => {
    render(<TripPlanner />);
    expect(screen.getByText(/Travel Itinerary Generator/i)).toBeInTheDocument();
  });

  test("selects a country and generates an itinerary", async () => {
    callApi.mockResolvedValue({
      status: 200,
      response: "Mocked itinerary response",
    });

    render(<TripPlanner />);

    fireEvent.change(screen.getByPlaceholderText("Thailand"), {
      target: { value: "Japan" },
    });
    fireEvent.change(screen.getByLabelText(/How many days?/i), {
      target: { value: 5 },
    });

    fireEvent.click(screen.getByText(/Generate/i));

    await waitFor(() => {
      expect(
        screen.getByText(/Mocked itinerary response/i)
      ).toBeInTheDocument();
    });
  });

  test("displays error when API call fails", async () => {
    callApi.mockResolvedValue({
      status: 400,
      error: "Mocked error message",
    });

    render(<TripPlanner />);

    fireEvent.click(screen.getByText(/Generate/i));

    await waitFor(() => {
      expect(screen.getByText(/Mocked error message/i)).toBeInTheDocument();
    });
  });

  test("redirects to home if token is invalid", async () => {
    callApi.mockResolvedValue({
      status: 400,
    });

    render(<TripPlanner />);

    await waitFor(() => {
      expect(mockPush).toHaveBeenCalledWith("/");
    });
  });
});
