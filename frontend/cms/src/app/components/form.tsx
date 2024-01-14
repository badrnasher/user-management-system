import React, { useEffect, useState } from "react";
import axios from "axios";

export default function Form({
  setShowForm,
  selectedUser,
  isEdit,
  fetchUsers,
}: {
  setShowForm: Function;
  selectedUser: any;
  isEdit: boolean;
  fetchUsers: Function;
}) {
  const [userName, setUserName] = useState("");
  const [userEmail, setUserEmail] = useState("");

  useEffect(() => {
    if (isEdit) {
      setUserName(selectedUser?.name || "");
      setUserEmail(selectedUser?.email || "");
    } else {
      setUserName("");
      setUserEmail("");
    }
  }, [isEdit]);

  const handleSubmit = async () => {
    try {
      if (isEdit) {
        // Update an existing user
        await axios.put(`/api/users/${selectedUser.id}`, {
          name: userName,
          email: userEmail,
        });
      } else {
        // Add a new user
        await axios.post(`/api/users`, {
          name: userName,
          email: userEmail,
        });
      }
      setShowForm(false); // Close the form after successful submission
      fetchUsers();
    } catch (error) {
      console.error("Error submitting user:", error);
    }
  };

  return (
    <>
      <form className="max-w-sm mx-auto bg-zinc-900 rounded-md p-8">
        <h1 className="text-center font-bold text-md mb-4">
          {isEdit ? "Edit User Form" : "Add New User Form"}
        </h1>
        <div className="mb-5">
          <label
            htmlFor="name"
            className="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
          >
            Name
          </label>
          <input
            value={userName}
            onChange={(e) => setUserName(e.target.value)}
            type="text"
            id="name"
            className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
            required
          />
        </div>
        <div className="mb-5">
          <label
            htmlFor="email"
            className="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
          >
            Email
          </label>
          <input
            value={userEmail}
            onChange={(e) => setUserEmail(e.target.value)}
            type="email"
            id="email"
            className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
            placeholder="name@email.com"
            required
          />
        </div>

        <div className="flex justify-between">
          <button
            onClick={handleSubmit}
            type="button"
            className={`text-white bg-blue-700 hover:bg-blue-800 ${
              isEdit
                ? "focus:ring-4 focus:outline-none focus:ring-blue-300"
                : ""
            } font-medium rounded-lg text-sm w-full sm:w-auto px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800 self-end`}
          >
            {isEdit ? "Update" : "Submit"}
          </button>
        </div>
      </form>
    </>
  );
}
