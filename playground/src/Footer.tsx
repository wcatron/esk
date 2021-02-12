import React, { useState } from "react";

import { HostSelector } from "./HostSelector";
import { Client } from 'esk-client-typescript';

export const Footer: React.FC<{
  host: string;
  client: Client;
  setHost: (value: string) => void;
  connected: boolean;
  subscribe: (topic: string) => void;
}> = ({ host, client, setHost, connected, subscribe }) => {
  const [topic, setTopic] = useState("");

  return (
    <footer className="App-footer">
      <HostSelector
        defaultValue={host}
        onSelect={setHost}
        connected={connected}
      />
      <input
        value={topic}
        placeholder={`Topic (ClientID: ${client && client.clientId})`}
        onChange={(e) => {
          setTopic(e.currentTarget.value);
        }}
      />
      <button
        onClick={() => {
          subscribe(topic);
          setTopic("");
        }}
      >
        Subscribe
      </button>
    </footer>
  );
};
