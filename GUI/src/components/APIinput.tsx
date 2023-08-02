import React, { useState } from "react";

class Error {
  error!: string;
}

class Ok {
  ok!: string;
}

type Resp = Ok | Error;

function IsError(r: Resp): r is Error {
  return (r as Error).error !== undefined;
}

const APIinput = (p: { url: string; endpoint: string; disabled: boolean }) => {
  const [value, setValue] = useState("");
  const [, setResult] = useState("");

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const response = await fetch(`${p.url}/${p.endpoint}?value=${value}`, {
      method: "POST",
    });
    const output = (await response.json()) as Resp;
    // let output = JSON.parse(data) as Resp;
    if (IsError(output)) {
      console.log(output.error);
    } else {
      console.log(output);
      setResult(output.ok);
    }
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
