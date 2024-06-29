"use client";
import React, { useEffect, useState } from "react";
import { FaUserCircle } from "react-icons/fa";
import Link from "next/link";

const NavBar = ({ color }) => {
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const [isAdmin, setIsAdmin] = useState(false);

  useEffect(() => {
    setIsAdmin(localStorage.getItem("type") == "Admin" ? true : false);
  }, []);

  const toggleDropdown = () => setIsDropdownOpen(!isDropdownOpen);

  return (
    <>
      {!color ? (
        <nav
          className={`bg-gray-800 text-white p-4 flex justify-between items-center sticky top-0`}
        >
          <Link href="/">
            <div className="text-lg font-semibold">
              Travel Itinerary Generator
            </div>
          </Link>
          <div className="relative">
            <button
              onClick={toggleDropdown}
              //   onMouseEnter={() => setIsDropdownOpen(true)}
              className="flex items-center focus:outline-none"
            >
              <FaUserCircle className="text-3xl" />
            </button>
            {isDropdownOpen && (
              <div
                className="absolute right-0 mt-2 py-2 w-48 bg-white rounded-md shadow-xl z-50 text-gray-800"
                onMouseLeave={() => setIsDropdownOpen(false)}
              >
                <Link
                  href="/dashboard"
                  className="block px-4 py-2 text-sm hover:bg-gray-100"
                >
                  Dashboard
                </Link>
                {isAdmin && (
                  <>
                    <Link
                      href="/user/admin"
                      className="block px-4 py-2 text-sm hover:bg-gray-100"
                    >
                      Manage
                    </Link>
                  </>
                )}
                {!isAdmin && (
                  <>
                    <Link
                      href="/contact"
                      className="block px-4 py-2 text-sm hover:bg-gray-100"
                    >
                      Support
                    </Link>
                  </>
                )}
                <Link
                  href="/"
                  className="block px-4 py-2 text-sm hover:bg-gray-100"
                  // onClick={() => setIsDropdownOpen(false)}
                  onClick={() => {
                    localStorage.clear();
                    // window.location.reload();
                  }}
                >
                  Logout
                </Link>
              </div>
            )}
          </div>
        </nav>
      ) : (
        <nav
          className={`bg-black text-white p-4 flex justify-between items-center sticky top-0`}
        >
          <Link href="/">
            <h1 className="text-xl font-bold">Travel Itinerary Generator</h1>
          </Link>
          <div
            className="relative"
            //   onMouseLeave={() => setIsDropdownOpen(false)}
          >
            <button
              onClick={toggleDropdown}
              //   onMouseEnter={() => setIsDropdownOpen(true)}
              className="flex items-center focus:outline-none"
            >
              <FaUserCircle className="text-3xl" />
            </button>
            {isDropdownOpen && (
              <div className="absolute right-0 mt-2 py-2 w-48 bg-white rounded-md shadow-xl z-50 text-gray-800">
                <Link
                  href="/dashboard"
                  className="block px-4 py-2 text-sm hover:bg-gray-100"
                >
                  Dashboard
                </Link>
                {isAdmin && (
                  <>
                    <Link
                      href="/user/admin"
                      className="block px-4 py-2 text-sm hover:bg-gray-100"
                    >
                      Manage
                    </Link>
                  </>
                )}
                {!isAdmin && (
                  <>
                    <Link
                      href="/contact"
                      className="block px-4 py-2 text-sm hover:bg-gray-100"
                    >
                      Support
                    </Link>
                  </>
                )}
                <Link
                  href="/"
                  className="block px-4 py-2 text-sm hover:bg-gray-100"
                  onClick={() => {
                    localStorage.clear();
                    window.location.reload();
                  }}
                >
                  Logout
                </Link>
              </div>
            )}
          </div>
        </nav>
      )}
    </>
  );
};

export default NavBar;
