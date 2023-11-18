import { useEffect, useState } from 'react';
import Markdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { PlaylistInfo, RawPlaylistInfo } from '../types';
import { BoxPlaylistInfo } from './BoxPlaylistInfo';

type FlashpointMarkdownProps = {
  children: string;
};

/**
 * Custom renderer for markdown that adds Flashpoint Community specific elements
 */
export function FlashpointMarkdown({ children }: FlashpointMarkdownProps) {
  return (
    <Markdown
      remarkPlugins={[remarkGfm]}
      components={{
        a: FlashpointLinkPreview
      }}>
      {children}
    </Markdown>
  );
}

// @TODO type safe this
/**
 * A custom link previewer for Flashpoint Community links
 */
export function FlashpointLinkPreview({ href, children, ...args }: any) {
  const [component, setComponent] = useState<JSX.Element>((
    <a href={href} {...args}>{children}</a>
  ));

  useEffect(() => {
    if (href.startsWith('/playlist/') && href.length > 10) {
      const idStr = href.slice(10);
      const id = parseInt(idStr);
      console.log(id);
      // Fetch playlist from the server
      fetch(`/api/playlist/${id}/preview`)
      .then((response) => response.json())
      .then((data) => {
        const rawPlaylist = data as RawPlaylistInfo;
        const playlist: PlaylistInfo = {
          id: rawPlaylist.id,
          name: rawPlaylist.name,
          description: rawPlaylist.description,
          totalGames: rawPlaylist.total_games,
          library: rawPlaylist.library,
          author: {
            authed: false,
            id: rawPlaylist.author.uid,
            username: rawPlaylist.author.username,
            avatarUrl: rawPlaylist.author.avatar_url,
            roles: rawPlaylist.author.roles,
            perms: [],
          }
        };
        return playlist;
      })
      .then((playlist) => {
        setComponent(playlist ? (
          <BoxPlaylistInfo playlist={playlist} />
        ) : (
          <h2>Failed to load playlist</h2>
        ));
      })
      .catch(() => {
        setComponent(<a href={href} {...args}>{children}</a>);
      });
    }
  }, [href]);

  return component;
}
