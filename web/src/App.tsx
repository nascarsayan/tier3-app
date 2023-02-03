import { useEffect, useState } from "preact/hooks";
import { JSXInternal } from "preact/src/jsx";

const BackendURL = process.env.REACT_APP_BACKEND_URL || "http://localhost:8081";

export function App() {
  const [fruits, setFruits] = useState({} as { [key: string]: number });

  useEffect(() => {
    fetch(BackendURL, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    })
      .then((response) => response.json())
      .then((data) => {
        console.log(data);
        setFruits(data["fruits"]);
      })
      .catch((error) => {
        console.error("Error:", error);
      });
  }, []);

  async function transact(
    e: JSXInternal.TargetedEvent<HTMLFormElement, Event>,
    mode: "buy" | "sell"
  ) {
    e.preventDefault();
    const form = e.target as HTMLFormElement;
    const formData = new FormData(form);
    const fruit = formData.get("fruit") as string;
    const quantity = formData.get("quantity") as string;
    try {
      let response = await fetch(`${BackendURL}/${mode}`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          fruit,
          quantity: parseInt(quantity),
        }),
      });
      if (!response.ok) {
        let msg = "HTTP status code: " + response.status;
        let t = await response.text();
        if (t) {
          msg += ", " + t;
        }
        let err = new Error(msg);
        throw err;
      }
      let data = await response.json();
      console.log(data);
      setFruits(data["fruits"]);
    } catch (error) {
      alert(error);
    }
  }

  function getTransactForm(mode: "buy" | "sell") {
    return (
      <form
        onSubmit={async (e) => {
          await transact(e, mode);
        }}
      >
        <label>
          Fruit
          <select name="fruit">
            <option value="apple">Apple</option>
            <option value="orange">Orange</option>
            <option value="banana">Banana</option>
          </select>
        </label>
        <label>
          Quantity
          <input type="number" name="quantity" />
        </label>
        <button type="submit">{mode}</button>
      </form>
    );
  }

  return (
    <div className="App">
      <section>
        {/* Section to buy fruits */}
        <h2>Buy fruits</h2>
        {getTransactForm("buy")}
      </section>
      <section>
        {/* Section to buy fruits */}
        <h2>Sell fruits</h2>
        {getTransactForm("sell")}
      </section>
      <section>
        <h2>Fruits in the inventory</h2>
        <ul>
          {Object.keys(fruits).map((fruit, idx) => (
            <li key={idx}>
              {fruit}: {fruits[fruit]}
            </li>
          ))}
        </ul>
      </section>
    </div>
  );
}
