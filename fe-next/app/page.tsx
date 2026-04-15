"use client";

export default function Home() {
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
