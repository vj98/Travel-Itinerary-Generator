"use client";
import NavBar from "@/components/navbar";
import React, { useEffect, useState } from "react";
import { callApi } from "@/app/utils/api-utils";
import { useRouter } from "next/navigation";
import toast, { Toaster } from "react-hot-toast";

const countries = ["Japan", "Italy", "France", "Spain", "Thailand"];
const months = [
  "January",
  "February",
  "March",
  "April",
  "May",
  "June",
  "July",
  "August",
  "September",
  "October",
  "November",
  "December",
];

const TripPlanner = () => {
  const [selectedCountry, setSelectedCountry] = useState("Thailand");
  const [selectedMonth, setSelectedMonth] = useState("April");
  const [days, setDays] = useState(7);
  const [recommendations, setRecommendations] = useState({
    restaurants: true,
    hotels: true,
  });
  const [content, setContent] = useState(null);
  const [loader, setLoader] = useState(false);
  const [token, setToken] = useState(null);
  const router = useRouter();

  const toggleRecommendation = (type) => {
    setRecommendations((prev) => ({ ...prev, [type]: !prev[type] }));
  };

  useEffect(() => {
    // comment this also for testing
    setToken(localStorage ? localStorage?.getItem("access_token") : null);
    try {
      async function ValidateToken() {
        const response = await callApi(
          process.env.NEXT_PUBLIC_TOKEN_VALIDATE,
          {},
          "GET"
        );
        if (response) {
          if (response.status != 200) {
            localStorage.clear();
            router.push("/");
          }
        }
      }
      if (localStorage?.getItem("access_token")) {
        ValidateToken();
      } else {
        router.push("/");
      }
    } catch (error) {
      console.log("error in checking auth ", error);
    }
  }, []);

  const handleGenerate = async () => {
    try {
      setLoader(true);

      let prompt = `Write me an itinerary for ${days} days 
    to ${selectedCountry} in the coming ${selectedMonth}. Describe the 
    weather that month, and also 5 things to take note 
    about this country's culture. Keep to a maximum travel area 
    to the size of Hokkaido, if possible, to minimize traveling time 
    between cities.\\n\\nFor each day, list me the following:\\n- Attractions 
    suitable for that season`;

      prompt += recommendations.hotels
        ? "\\n- Hotel (prefer not to change it unless traveling to another city)"
        : "";
      prompt += recommendations.restaurants
        ? `\\n- 2 Restaurants, one for lunch and another for dinner, with 
      shortened Google Map links\\nand give me 
    a daily summary of the above points into a paragraph or two.`
        : "";

      const body = JSON.stringify({
        prompt: prompt,
        session: localStorage.getItem("session"),
      });
      const payload = {
        body: body,
      };

      const response = await callApi(
        process.env.NEXT_PUBLIC_FETCH_DATA_OPEN_AI,
        payload,
        "POST"
      );
      console.log("response ", response);
      if (response.status != 200) {
        toast.error(response.error);
      } else {
        setContent(`${response.response}`);
      }
    } catch (err) {
      console.log("error in calling generate ", err);
      toast.error(err.message);
    } finally {
      setLoader(false);
    }
  };

  const handleCountry = (e) => {
    setContent("");
    setSelectedCountry(e);
  };

  const handleChangeMonth = (e) => {
    setSelectedMonth(e.target.value);
  };

  return (
    <>
      {token && ( // comment this also for testing
        <>
          <NavBar /> {/*  for test case comment this line*/}
          <Toaster position="top-center" reverseOrder={false} />
          <div className="flex lg:flex-row flex-col h-screen bg-[#1A1A2E] text-white">
            <div className="flex flex-col lg:w-2/5 w-100 p-10 bg-[#16213E] ">
              <h1 className="text-4xl text-green-500 mb-4">
                Travel Itinerary Generator
              </h1>
              <p className="mb-8 text-blue-300">
                {
                  "Give me some details and I'll craft an itinerary just for you!"
                }
              </p>

              <div className="mb-6">
                <label className="mb-2 block">
                  1. Where do you want to go?
                </label>
                <input
                  type="text"
                  placeholder="Thailand"
                  className="w-full p-2 mb-4 bg-[#1A1A2E] rounded-md"
                  value={selectedCountry}
                  onChange={(e) => handleCountry(e.target.value)}
                />
                <div className="flex justify-between flex-wrap">
                  {countries.map((country) => (
                    <button
                      key={country}
                      className={`mt-2 px-4 py-2 border border-green-500 rounded-md ${
                        selectedCountry === country ? "bg-green-500" : ""
                      }`}
                      onClick={() => handleCountry(country)}
                    >
                      {country}
                    </button>
                  ))}
                </div>
              </div>

              <div className="mb-6">
                <label className="mb-2 block" for="How many days?">
                  2. How many days?
                </label>
                <input
                  type="number"
                  className="w-full p-2 mb-4 bg-[#1A1A2E] rounded-md"
                  value={days}
                  onChange={(e) => setDays(e.target.value)}
                />
              </div>

              <div className="mb-6">
                <label className="mb-2 block">3. Month</label>
                <select
                  className="w-full p-2 bg-[#1A1A2E] rounded-md"
                  value={selectedMonth}
                  onChange={(e) => handleChangeMonth(e)}
                >
                  {months.map((month) => (
                    <option key={month} value={month}>
                      {month}
                    </option>
                  ))}
                </select>
              </div>

              <div className="mb-6">
                <label className="mb-2 block">4. Recommendations?</label>
                <div className="flex gap-4  flex-wrap">
                  <button
                    className={`px-6 py-2 border border-blue-500 rounded-md ${
                      recommendations.restaurants ? "bg-blue-500" : ""
                    }`}
                    onClick={() => toggleRecommendation("restaurants")}
                  >
                    Restaurants
                  </button>
                  <button
                    className={`px-6 py-2 border border-blue-500 rounded-md ${
                      recommendations.hotels ? "bg-blue-500" : ""
                    }`}
                    onClick={() => toggleRecommendation("hotels")}
                  >
                    Hotels
                  </button>
                </div>
              </div>

              {!loader ? (
                <button
                  className="w-full py-3 bg-red-600 rounded-md hover:bg-red-700"
                  onClick={handleGenerate}
                >
                  Generate
                </button>
              ) : (
                <div role="status" className="flex justify-center">
                  <svg
                    aria-hidden="true"
                    className="w-8 h-8 text-gray-200 animate-spin dark:text-gray-600 fill-red-600"
                    viewBox="0 0 100 101"
                    fill="none"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path
                      d="M100 50.5908C100 78.2051 77.6142 100.591 50 100.591C22.3858 100.591 0 78.2051 0 50.5908C0 22.9766 22.3858 0.59082 50 0.59082C77.6142 0.59082 100 22.9766 100 50.5908ZM9.08144 50.5908C9.08144 73.1895 27.4013 91.5094 50 91.5094C72.5987 91.5094 90.9186 73.1895 90.9186 50.5908C90.9186 27.9921 72.5987 9.67226 50 9.67226C27.4013 9.67226 9.08144 27.9921 9.08144 50.5908Z"
                      fill="currentColor"
                    />
                    <path
                      d="M93.9676 39.0409C96.393 38.4038 97.8624 35.9116 97.0079 33.5539C95.2932 28.8227 92.871 24.3692 89.8167 20.348C85.8452 15.1192 80.8826 10.7238 75.2124 7.41289C69.5422 4.10194 63.2754 1.94025 56.7698 1.05124C51.7666 0.367541 46.6976 0.446843 41.7345 1.27873C39.2613 1.69328 37.813 4.19778 38.4501 6.62326C39.0873 9.04874 41.5694 10.4717 44.0505 10.1071C47.8511 9.54855 51.7191 9.52689 55.5402 10.0491C60.8642 10.7766 65.9928 12.5457 70.6331 15.2552C75.2735 17.9648 79.3347 21.5619 82.5849 25.841C84.9175 28.9121 86.7997 32.2913 88.1811 35.8758C89.083 38.2158 91.5421 39.6781 93.9676 39.0409Z"
                      fill="currentFill"
                    />
                  </svg>
                  <span className="sr-only">Loading...</span>
                </div>
              )}
            </div>
            {/* <div className="w-3/5 p-10 bg-[#1A1A2E]">
        {/* The right panel can be used for displaying the generated itinerary */}
            {/* </div> */}
            <RightPanel
              country={selectedCountry}
              days={days}
              content={content}
            />
          </div>
        </>
      )}
      {/*comment this also for testing*/}
    </>
  );
};

const RightPanel = ({ country, days, content }) => {
  const itineraryLines = content?.split("\n").map((line, index) => (
    <React.Fragment key={index}>
      {line}
      <br />
    </React.Fragment>
  ));
  return (
    <div className="lg:w-3/5 w-100 p-10 mb-16 text-white lg:overflow-auto bg-[#1A1A2E]">
      <div className="scroll-smooth">
        {content && (
          <>
            <h2 className="text-2xl font-semibold mb-4">
              Your iteneray of {country} for {days} days
            </h2>
            <span>{itineraryLines}</span>
          </>
        )}
      </div>
    </div>
  );
};

export default TripPlanner;
