"use client";
import { callApi } from "@/app/utils/api-utils";
import NavBar from "@/components/navbar";
import React, { useEffect, useState } from "react";

const Admin = () => {
  const [isAdmin, setIsAdmin] = useState(false);
  const [allUsers, setAllUsers] = useState([]);
  const [status, setStatus] = useState("");
  const [loaderStatus, setLoaderStatus] = useState(false);

  useEffect(() => {
    setIsAdmin(localStorage?.getItem("type") == "Admin" ? true : false);
    handleFetchUser();
  }, []);

  const handleStatus = async (e, val) => {
    e.preventDefault();
    try {
      setLoaderStatus(true);
      setStatus(e.target.value);
      const body = JSON.stringify({
        email: val.Email,
        status: e.target.value,
        session: localStorage.getItem("session"),
      });
      const payload = {
        body: body,
      };

      console.log(body, payload);
      const resp = await callApi(
        process.env.NEXT_PUBLIC_UPDATE_STATUS,
        payload,
        "POST"
      );
      console.log("updated response", resp);
    } catch (err) {
      console.log("error in status update ", err);
    } finally {
      setLoaderStatus(false);
      handleFetchUser();
    }
  };

  const handleFetchUser = async () => {
    try {
      const body = JSON.stringify({
        session: localStorage.getItem("session"),
      });
      const payload = {
        body: body,
      };
      const response = await callApi(
        process.env.NEXT_PUBLIC_FETCH_ALL_USERS,
        payload,
        "POST"
      );
      if (response) {
        if (response.status == "200") {
          setAllUsers(response.data);
        }
      }
    } catch (error) {
      console.log("error while fetching all users ", error);
    }
  };
  return (
    <>
      <NavBar />
      {isAdmin && (
        <div className="mx-auto max-w-screen-lg px-4 py-8 sm:px-8">
          <div className="flex items-center justify-between pb-6">
            <div>
              <h2 className="font-semibold text-gray-700">User Accounts</h2>
              <span className="text-xs text-gray-500">
                View accounts of registered users
              </span>
            </div>
            {/* <div className="flex items-center justify-between">
          <div className="ml-10 space-x-8 lg:ml-40">
            <button className="flex items-center gap-2 rounded-md bg-blue-600 px-4 py-2 text-sm font-semibold text-white focus:outline-none focus:ring hover:bg-blue-700">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke-width="1.5"
                stroke="currentColor"
                className="h-4 w-4"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M12 4.5v15m0 0l6.75-6.75M12 19.5l-6.75-6.75"
                />
              </svg>
              CSV
            </button>
          </div>
        </div> */}
          </div>
          <div className="overflow-y-hidden rounded-lg border">
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="bg-blue-600 text-left text-xs font-semibold uppercase tracking-widest text-white">
                    <th className="px-5 py-3">ID</th>
                    <th className="px-5 py-3">Email ID</th>
                    <th className="px-5 py-3">User Role</th>
                    <th className="px-5 py-3">Created at</th>
                    <th className="px-5 py-3">Status</th>
                  </tr>
                </thead>
                <tbody className="text-gray-500">
                  {allUsers.map((val, i) => {
                    return (
                      <tr key={val.ID}>
                        <td className="border-b border-gray-200 bg-white px-5 py-5 text-sm">
                          <p className="whitespace-no-wrap">{i}</p>
                        </td>
                        <td className="border-b border-gray-200 bg-white px-5 py-5 text-sm">
                          <div className="flex items-center">
                            <div className="ml-3">
                              <p className="whitespace-no-wrap">{val.Email}</p>
                            </div>
                          </div>
                        </td>
                        <td className="border-b border-gray-200 bg-white px-5 py-5 text-sm">
                          <p className="whitespace-no-wrap">{val.Type}</p>
                        </td>
                        <td className="border-b border-gray-200 bg-white px-5 py-5 text-sm">
                          <p className="whitespace-no-wrap">
                            {val.LastUpdatedAt}
                          </p>
                        </td>

                        <td className="border-b border-gray-200 bg-white px-5 py-5 text-sm">
                          <select
                            className={`${
                              val.Type == "Admin" ? "" : "cursor-pointer"
                            } rounded-full px-3 py-1 text-sm font-semibold ${
                              val.Status == "Pending"
                                ? "text-orange-400"
                                : val.Status == "InActive"
                                ? "text-red-400"
                                : "text-green-900"
                            }`}
                            value={val.Status}
                            onChange={(e) => handleStatus(e, val)}
                            disabled={loaderStatus || val.Type == "Admin"}
                          >
                            <option value="Pending">Pending</option>
                            <option value="Active">Active</option>
                            <option value="InActive">InActive</option>
                          </select>
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
            {/* <div className="flex flex-col items-center border-t bg-white px-5 py-5 sm:flex-row sm:justify-between">
            <span className="text-xs text-gray-600 sm:text-sm">
              {" "}
              Showing 1 to 5 of 12 Entries{" "}
            </span>
            <div className="mt-2 inline-flex sm:mt-0">
              <button className="mr-2 h-12 w-12 rounded-full border text-sm font-semibold text-gray-600 transition duration-150 hover:bg-gray-100">
                Prev
              </button>
              <button className="h-12 w-12 rounded-full border text-sm font-semibold text-gray-600 transition duration-150 hover:bg-gray-100">
                Next
              </button>
            </div>
          </div> */}
          </div>
        </div>
      )}
    </>
  );
};

export default Admin;
