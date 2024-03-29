// import PageOne from "./PageOne"
import { useState } from "react";
import APIinput from "./components/APIinput";
import APIButton from "./components/APiButton";
import ChartContainerZ from "./components/Chart";
import StatusBar from "./components/Status";

const url = "http://localhost:8080";
const ws = "ws://localhost:8080/ws";

function Page() {
  const [isStarted, setIsStarted] = useState(false);
  const setIsStopped = () => setIsStarted(false);
  const [error, setError] = useState<string | undefined>();

  return (
    <div className="container">
      <div className="header">
        <APIButton
          url={url}
          endpoint="start"
          disabled={isStarted}
          setValue={setIsStarted}
          setError={setError}
        />
        <APIButton
          url={url}
          endpoint="stop"
          disabled={!isStarted}
          setValue={setIsStopped}
          setError={setError}
        />
        <APIinput endpoint="temperature" disabled={!isStarted} url={url} />
        <APIinput endpoint="pressure" disabled={!isStarted} url={url} />
      </div>

      <div className="mainContent">
        {isStarted && <ChartContainerZ URL={ws} />}
      </div>
      <div className="footer">
        <StatusBar URL={ws} started={isStarted} />
      </div>
    </div>
  );
}

export default Page;
