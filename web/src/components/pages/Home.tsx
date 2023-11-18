import { Button, Grid, Icon, Loader } from 'semantic-ui-react';
import { NewsPostBox } from '../NewsPost';
import { SideBar } from '../SideBar';
import { PermissionLock } from '../Locks';
import { UserPerm } from '../../utils/permissions';
import { NewsPost } from '../../types';
import { useEffect, useState } from 'react';
import { ResponsePosts } from '../../resTypes';
import { mapRawNewsPost } from '../../utils/mappers';
import { useSmartNavigate } from '../../hooks/useSmartNavigate';

export function Home() {
  const [loadingPosts, setLoadingPosts] = useState<boolean>(true);
  const [posts, setPosts] = useState<NewsPost[]>([]);
  const navigate = useSmartNavigate();

  useEffect(() => {
    // Fetch up to 5 posts
    fetch('/api/posts?page=1&page_size=5&order_by=updated_at&order_direction=desc')
    .then(res => res.json())
    .then((data: ResponsePosts) => {
      setPosts(data.posts.map(mapRawNewsPost));
      setLoadingPosts(false);
    });
  },[]);

  const gotoLink = (link: string) => {
    window.open(link, '_blank');
  };

  return (
    <div>
      <Grid>
        <Grid.Column width={3}>
          <SideBar title={'Routes'}>
            <div className='sidebar__item'>
              <p><b>Need help? Visit our Discord</b></p>
              <Button onClick={() => gotoLink('https://discord.com/invite/qhvAkhWXU5')} color='linkedin' fluid>
                <Icon name='discord' /> Discord
              </Button>
            </div>
            <div className='sidebar__item'>
              <Button onClick={() => gotoLink('https://flashpointarchive.org/datahub/Main_Page')} fluid>
                <Icon name='wikipedia w' /> Datahub
              </Button>
            </div>
          </SideBar>
          <PermissionLock perm={UserPerm.CREATE_NEWS_POST}>
            <SideBar title={'Tools'}>
              <Button onClick={() => navigate('/post/create')} fluid>
                <Icon name='newspaper' /> Create News Post
              </Button>
            </SideBar>
          </PermissionLock>
        </Grid.Column>
        <Grid.Column width={13}>
          {
            !loadingPosts ?
              posts.length > 0 ? (
                posts.map((post, index) => (
                  <NewsPostBox key={index} post={post} />
                ))
              ) : (
                <h2>No news posts available</h2>
              ) : (
                <Loader active inline='centered'>Fetching Posts...</Loader>
              )
          }
        </Grid.Column>
      </Grid>
    </div>
  );
}
