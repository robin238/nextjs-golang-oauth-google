import { useEffect, useState } from "react";
import axios from "axios";

export default function Dashboard() {
  const [data, setData] = useState("");

  useEffect(() => {
    axios
      .get("http://localhost:8080/dashboard")
      .then((res) => setData(res.data.message));
  }, []);

  return (
    <div>
      <h1>Dashboard</h1>
      <p>{data}</p>
    </div>
  );
}
