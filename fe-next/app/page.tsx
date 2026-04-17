"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";

export default function Home() {
  const router = useRouter();

  useEffect(() => {
    const token = localStorage.getItem("authToken");
    if (token) {
      router.replace("/dashboard");
    }
  }, [router]);

  const loginGoogle = () => {
    window.location.href = "http://localhost:8080/auth/google";
  };

  return (
    <div style={{ padding: 40 }}>
      <h1>Login / Register</h1>
      <button onClick={loginGoogle} className="bg-blue-300 p-2 border-2">
        Login dengan Google
      </button>
    </div>
  );
}
