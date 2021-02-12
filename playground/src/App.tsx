import React, { useEffect, useState } from "react";
import "./App.css";
import { Client } from "esk-client-typescript";
import { TopicCard } from "./TopicCard";
import { Footer } from "./Footer";
let client: Client;

const useEsk = (defaultHost: string) => {
  const [host, setHost] = useState(defaultHost);
  const [connected, setConnected] = useState(false);
  const [topics, setTopics] = useState<string[]>([]);

  useEffect(() => {
    console.log("Connecting...");
    client = new Client({
      url: "ws://" + host + "/ws",
    });
    client.on("open", () => {
      setConnected(true);
      console.log("Connected");
    });
    client.on("close", () => {
      setConnected(false);
    });
    return () => {
      client.disconnect();
    };
  }, [host]);

  const subscribe = (topic: string) => {
    setTopics([...topics, topic]);
  };

  return { host, setHost, connected, topics, client, subscribe };
};

function App() {
  const { host, setHost, connected, topics, subscribe } = useEsk(
    window.location.host
  );

  return (
    <div
      className="App"
      style={{
        background: connected ? "#282c34" : "#484c54",
      }}
    >
      <div>
        {topics.map((topic) => (
          <TopicCard
            key={topic}
            topic={topic}
            client={client}
            onPublish={(payload: string) => {
              client.publish(topic, payload);
            }}
          />
        ))}
      </div>
      <Footer
        host={host}
        setHost={setHost}
        client={client}
        connected={connected}
        subscribe={subscribe}
      />
    </div>
  );
}

export default App;
