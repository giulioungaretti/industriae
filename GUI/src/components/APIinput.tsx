import React, { useState } from "react";

const APIinput = (p: { url: string; endpoint: string; disabled: boolean }) => {
  const [value, setValue] = useState("");
  const [, setResult] = useState("");

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const response = await fetch(`${p.url}/${p.endpoint}?value=${value}`, {
      method: "POST",
    });
    const data = await response.json();
    setResult(data);
  };

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setValue(event.target.value);
  };
  return (
    <div className="input-box">
      <form onSubmit={handleSubmit}>
        <label>
          {p.endpoint}:
          <input type="number" value={value} onChange={handleChange} />
        </label>
        <button type="submit" disabled={p.disabled}>
          Set
        </button>
      </form>
    </div>
  );
};

export default APIinput;
