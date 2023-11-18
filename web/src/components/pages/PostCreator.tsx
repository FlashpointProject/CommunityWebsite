import { useState } from 'react';
import { Button, Form, FormProps } from 'semantic-ui-react';
import { FlashpointMarkdown } from '../FlashpointMarkdown';
import { SegmentWithHeader } from '../SegmentWithHeader';
import { NewsPost } from '../../types';
import { useSmartNavigate } from '../../hooks/useSmartNavigate';

export function PostCreator() {
  const [preview, setPreview] = useState(false);
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const navigate = useSmartNavigate();

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>, data: FormProps) => {
    fetch('/api/posts', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        title: title,
        content: content,
        post_type: 'news'
      })
    })
    .then((res) => {
      if (!res.ok) {
        throw new Error('Failed to create post');
      }
      return res;
    })
    .then((res) => res.json())
    .then((post: NewsPost) => {
      navigate(`/post/${post.id}`);
    });
  };

  return (
    <SegmentWithHeader header='Create News Post'>
      <Form onSubmit={handleSubmit}>
        <Form.Input
          label='Title'
          placeholder='Title'
          value={title}
          onChange={(e, data) => { setTitle(data.value as string); }}
          required
        />
        <Form.TextArea
          label='Content'
          placeholder='Content'
          value={content}
          onChange={(e, data) => { setContent(data.value as string); }}
          required />
        <div className='button-row'>
          <Button type="button" onClick={() => { setPreview(!preview); }}>Preview</Button>
          <Form.Button primary>Submit</Form.Button>
        </div>
        { preview && (
          <FlashpointMarkdown>{content}</FlashpointMarkdown>
        )}
      </Form>
    </SegmentWithHeader>
  );
}
