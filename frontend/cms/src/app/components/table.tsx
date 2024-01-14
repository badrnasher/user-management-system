import React, { useEffect, useState } from "react";
import Form from "./form";
import axios from "axios";

//set default axios base url
axios.defaults.baseURL = "http://localhost:8080";

export interface User {
  id: number;
  name: string;
  email: string;
}
export default function Table() {
  const [showForm, setShowForm] = useState(false);
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [isEdit, setIsEdit] = useState(false);
  const [users, setUsers] = useState<User[]>([]);

  const fetchUsers = async () => {
    try {
      const response = await axios.get("/api/users");
      setUsers(response.data);
    } catch (error) {
      console.error("Error fetching users:", error);
    }
  };

  useEffect(() => {
    // Fetch users from the API on component mount
    fetchUsers();
  }, []);

  const handleNewUser = () => {
    setShowForm(true);
    setIsEdit(false);
  };

  const handleEditUser = (user: User) => {
    setSelectedUser(user);
    setIsEdit(true);
    setShowForm(true);
  };

  const handleRemove = async (userID: Number) => {
    try {
      await axios.delete(`/api/users/${userID}`);
      setShowForm(false);
      fetchUsers();
    } catch (error) {
      console.error("Error deleting user:", error);
    }
  };

  return (
    <div className="relative overflow-x-auto shadow-md sm:rounded-lg w-2/3 mt-24 mx-auto">
      {!showForm ? (
        <>
          <div className="flex justify-end my-2">
            <button
              className="px-2 py-1 bg-indigo-500 rounded-md "
              onClick={handleNewUser}
            >
              Add new user
            </button>
          </div>
          <table className="w-full text-sm text-left rtl:text-right text-gray-500 dark:text-gray-400">
            <thead className="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
              <tr>
                <th scope="col" className="px-6 py-3">
                  ID
                </th>
                <th scope="col" className="px-6 py-3">
                  Name
                </th>
                <th scope="col" className="px-6 py-3">
                  Email
                </th>
                <th scope="col" className="px-6 py-3">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody>
              {users.map((user) => (
                <tr
                  key={user.id}
                  className="bg-white border-b dark:bg-gray-800 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600"
                >
                  <td className="px-6 py-4">{user.id}</td>
                  <td className="px-6 py-4">{user.name}</td>
                  <td className="px-6 py-4">{user.email}</td>
                  <td className="flex items-center px-6 py-4">
                    <span
                      className="font-medium cursor-pointer text-blue-600 dark:text-blue-500 hover:underline"
                      onClick={() => handleEditUser(user)}
                    >
                      Edit
                    </span>
                    <span
                      className="font-medium cursor-pointer text-red-600 dark:text-red-500 hover:underline ms-3"
                      onClick={() => handleRemove(user.id)}
                    >
                      Delete
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </>
      ) : (
        <Form
          setShowForm={setShowForm}
          selectedUser={selectedUser}
          isEdit={isEdit}
          fetchUsers={fetchUsers}
        />
      )}
    </div>
  );
}
