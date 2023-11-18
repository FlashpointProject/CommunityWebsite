import { useCallback, useState } from 'react';
import Markdown from 'react-markdown';
import { Button, Dropdown, Grid, Input, Segment } from 'semantic-ui-react';

export function TestPage() {
  const [url, setUrl] = useState('https://fpcomm-dev.colintest.site/api/');
  const [body, setBody] = useState('');
  const [method, setMethod] = useState('GET');
  const [response, setResponse] = useState('no request sent yet');

  const onSubmit = useCallback(() => {
    const b = method === 'POST' && body !== '' ? body : undefined;
    fetch(url, {
      method,
      body: b,
      headers: {
        'Content-Type': 'application/json'
      }
    })
    .then(async (res) => {
      try {
        if (res.ok) {
          const data = await res.json();
          setResponse(JSON.stringify(data, null, 2));
        } else {
          setResponse(`Error ${res.status}: ${await res.text()}`);
        }
      } catch (e) {
        setResponse(`Error reading response: ${e}`);
      }
    })
    .catch((e) => {
      setResponse(`Error: ${e}`);
    });
  }, [url, body, method]);

  return (
    <Segment>
      <h2>Web Request Tester</h2>
      <Grid>
        <Grid.Row>
          <Grid.Column width={2}>
            <b>Request Type</b>
          </Grid.Column>
          <Grid.Column width={14}>
            <Dropdown
              placeholder="request type"
              value={method}
              onChange={(event, data) => {
                setMethod(data.value as string);
              }}
              options={[
                {key: 'GET', text: 'GET', value: 'GET'},
                {key: 'POST', text: 'POST', value: 'POST'}
              ]}
              selection/>
          </Grid.Column>
        </Grid.Row>
        <Grid.Row>
          <Grid.Column width={2}>
            <b>Request URL</b>
          </Grid.Column>
          <Grid.Column width={14}>
            <Input
              fluid
              placeholder="request url"
              onChange={(event, data) => { setUrl(data.value as string); }}
              value={url}/>
          </Grid.Column>
        </Grid.Row>
        <Grid.Row>
          <Grid.Column width={2}>
            <b>Request Body</b>
          </Grid.Column>
          <Grid.Column width={14}>
            <Input
              fluid
              placeholder="request body"
              onChange={(event, data) => { setBody(data.value as string); }}
              value={body}/>
          </Grid.Column>
        </Grid.Row>
        <Grid.Row>
          <Grid.Column>
            <Button onClick={onSubmit}>Submit</Button>
          </Grid.Column>
        </Grid.Row>
      </Grid>
      <div>
        <Markdown>{'```\n' + response + '\n```'}</Markdown>
      </div>
    </Segment>
  );
}
