import React, { useEffect, useState } from 'react';
import './App.css';
import { Client } from 'esk-client-typescript';
import { HostSelector } from './HostSelector';
import { TopicCard } from './TopicCard';

let client:Client

function App() {
  const [ connected, setConnected ] = useState(false)
  const [ topic, setTopic ] = useState('')
  const [ topics, setTopics ] = useState<string[]>([])
  const [ host, setHost ] = useState(window.location.host)
  useEffect(() => {
    console.log('Connecting...')
    client = new Client({
      url: 'ws://'+host+'/ws'
    })
    client.on('open', () => {
      setConnected(true)
      console.log('Connected')
    })
    client.on('close', () => {
      setConnected(false)
    })
    return () => {
      client.disconnect()
    }
  }, [host])

  const subscribe = (topic: string) => {
    setTopics([...topics, topic])
  }

  return (
    <div className="App" style={{
      background: connected ? '#282c34' : '#484c54'
    }}>
      <div>
        {topics.map(topic => <TopicCard key={topic} topic={topic} client={client} onPublish={(payload: string) => {
          client.publish(topic, payload)
        }} />)}
      </div>
      <footer className="App-footer">
        <HostSelector defaultValue={host} onSelect={setHost} connected={connected} />
        <input value={topic} placeholder={`Topic (ClientID: ${client && client.clientId})`} onChange={(e) => {
          setTopic(e.currentTarget.value)
        }} /><button onClick={() => {
          subscribe(topic)
          setTopic('')
        }}>Subscribe</button>
      </footer>
    </div>
  );
}

export default App;
