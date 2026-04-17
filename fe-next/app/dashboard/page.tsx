"use client";

import { useEffect, useState } from "react";
import axios from "axios";

export default function Dashboard() {
  const [userData, setUserData] = useState<{
    message: string;
    email?: string;
    name?: string;
    role?: string;
  }>({
    message: "Loading...",
  });

  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    const urlToken = params.get("token");
    const storageToken = localStorage.getItem("authToken");
    const token = urlToken || storageToken;

    if (!token) {
      setUserData({ message: "Anda belum login." });
      return;
    }

    if (urlToken) {
      localStorage.setItem("authToken", urlToken);
      params.delete("token");
      const cleanUrl = `${window.location.pathname}${params.toString() ? `?${params.toString()}` : ""}`;
      window.history.replaceState({}, "", cleanUrl);
    }

    axios
      .get("http://localhost:8080/dashboard", {
        headers: { Authorization: `Bearer ${token}` },
      })
      .then((res) => setUserData(res.data))
      .catch(() => {
        localStorage.removeItem("authToken");
        setUserData({ message: "Sesi tidak valid. Silakan login ulang." });
      });
  }, []);

  const handleLogout = async () => {
    const token = localStorage.getItem("authToken");

    if (token) {
      try {
        await axios.post("http://localhost:8080/logout", null, {
          headers: { Authorization: `Bearer ${token}` },
        });
      } catch (error) {
        // Continue logout even if backend logout fails
      }
    }

    localStorage.removeItem("authToken");
    window.location.href = "/";
  };

  return (
    <div style={{ padding: 40 }}>
      <h1>Dashboard</h1>
      <p>{userData.message}</p>
      {userData.email && <p>Email: {userData.email}</p>}
      {userData.name && <p>Name: {userData.name}</p>}
      {userData.role && <p>Role: {userData.role}</p>}
      <button onClick={handleLogout} className="bg-red-300 p-2 border-2 mt-4">
        Logout
      </button>
    </div>
  );
}
