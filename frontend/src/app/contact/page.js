"use client";
import NavBar from "@/components/navbar";
import React, { useEffect, useState } from "react";
import { callApi } from "../utils/api-utils";

const Contact = () => {
  const [email, setEmail] = useState("");
  const [subject, setSubject] = useState("");
  const [content, setContent] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSubj = (e) => {
    e.preventDefault();
    setSubject(e.target.value);
  };

  const handleContent = (e) => {
    e.preventDefault();
    setContent(e.target.value);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    const temp = content + "<br/> <b>Email orignally from " + email + "</b>";
    setLoading(true);
    try {
      const body = JSON.stringify({
        email: email,
        subject: subject,
        content: content,
        session: localStorage.getItem("session"),
      });
      const payload = {
        body: body,
      };
      const resp = await callApi(
        process.env.NEXT_PUBLIC_SEND_MAIL_CONTACT,
        payload,
        "POST"
      );

      console.log("resp sent mail ", resp);
      if (resp.status == 200) {
        setSubject("");
        setContent("");
      }
    } catch (err) {
      console.log("error in sending mail ", err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    // setEmail(localStorage.getItem("email"));
    setEmail(localStorage?.getItem("email"));
    // console.log(localStorage);
  }, []);

  return (
    <>
      <NavBar />
      <div className=" bg-gray-900" style={{ height: "100vh" }}>
        <div className="py-8 lg:py-16 px-4 mx-auto max-w-screen-md">
          <h2 className="mb-4 text-4xl tracking-tight font-extrabold text-center text-white ">
            Contact Us
          </h2>
          <p className="mb-8 lg:mb-16 font-light text-center text-gray-500 sm:text-xl">
            Got a technical issue? Want to send feedback about a beta feature?
            Need details about our Business plan? Let us know.
          </p>
          <form action="#" className="space-y-8">
            <div>
              <label
                for="email"
                className="block mb-2 text-sm font-medium text-gray-900 "
              >
                Your email
              </label>
              <input
                type="email"
                id="email"
                className="shadow-sm bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-primary-500 focus:border-primary-500 block w-full p-2.5 "
                placeholder="your@mail.com"
                required
                disabled
                value={email}
              />
            </div>
            <div>
              <label
                for="subject"
                className="block mb-2 text-sm font-medium text-gray-900 "
              >
                Subject
              </label>
              <input
                type="text"
                id="subject"
                className="block p-3 w-full text-sm text-gray-900 bg-gray-50 rounded-lg border border-gray-300 shadow-sm focus:ring-primary-500 focus:border-primary-500 "
                placeholder="Let us know how we can help you"
                required
                value={subject}
                onChange={handleSubj}
              />
            </div>
            <div className="sm:col-span-2">
              <label
                for="message"
                className="block mb-2 text-sm font-medium text-gray-900 "
              >
                Your message
              </label>
              <textarea
                id="message"
                rows="6"
                className="block p-2.5 w-full text-sm text-gray-900 bg-gray-50 rounded-lg shadow-sm border border-gray-300 focus:ring-primary-500 focus:border-primary-500 "
                placeholder="Leave a comment..."
                value={content}
                onChange={handleContent}
              ></textarea>
            </div>
            <button
              disabled={loading}
              onClick={handleSubmit}
              type="submit"
              className="border-white outline-dashed py-3 px-5 text-sm font-medium text-center text-white rounded-lg bg-primary-700 sm:w-fit hover:bg-primary-800 focus:ring-4 focus:outline-none focus:ring-primary-300 "
            >
              {loading ? "Sending..." : "Send message"}
            </button>
          </form>
        </div>
      </div>
    </>
  );
};

export default Contact;
