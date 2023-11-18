import { useEffect, useState } from 'react';
import { NewsPost, RawNewsPost } from '../../types';
import { NewsPostBox } from '../NewsPost';
import { NotFound } from './NotFound';
import { useParams } from 'react-router-dom';
import { mapRawNewsPost } from '../../utils/mappers';
import { Loader } from 'semantic-ui-react';

export function PostPage() {
  const { id } = useParams();
  const [loading, setLoading] = useState(true);
  const [post, setPost] = useState<NewsPost>();

  useEffect(() => {
    fetch(`/api/post/${id}`)
    .then(res => res.json())
    .then((data: RawNewsPost) => {
      setPost(mapRawNewsPost(data));
      setLoading(false);
    });
  }, []);

  return !loading ?
    post ? (
      <NewsPostBox post={post!} />
    ) : (
      <NotFound/>
    ) : (
      <Loader active inline='centered'>Loading Post...</Loader>
    );
}
