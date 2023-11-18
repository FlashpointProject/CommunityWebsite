import { Button, Item, Label } from 'semantic-ui-react';
import { Game, PlaylistGame } from '../types';
import { useMemo, useState } from 'react';

export type GameListProps = {
  games: Game[];
  expanded?: boolean;
};

export type GameListPlaylistProps = {
  games: PlaylistGame[];
  expanded?: boolean;
};

/**
 * Itemized list of games (sister to GameList)
 * @param games Playlist Games to display
 */
export function GameListPlaylist(props: GameListPlaylistProps) {
  const [expanded, setExpanded] = useState(props.expanded || false);

  const games = useMemo(() => {
    if (expanded || props.games.length <= 5) {
      return props.games.map((game) => formatGame(game.game, game.notes));
    } else {
      return [...props.games.slice(0, 5).map((game) => formatGame(game.game, game.notes)),
        (
          <Button size='huge' onClick={() => setExpanded(true)} fluid>View {props.games.length - 5} more games...</Button>
        )];
    }
  }, [expanded, props.games]);

  if (props.games.length === 0) {
    return <p>No games present</p>;
  } else {
    return (
      <Item.Group divided>
        {games}
      </Item.Group>
    );
  }
}

/**
 * Itemized list of games (sister to GameListPlaylist)
 * @param games Games to display
 */
export function GameList(props: GameListProps) {
  const [expanded, setExpanded] = useState(props.expanded || false);

  const games = useMemo(() => {
    if (expanded || props.games.length <= 5) {
      return props.games.map((game) => formatGame(game));
    } else {
      return [...props.games.slice(0, 5).map((game) => formatGame(game)),
        (
          <Button size='huge' onClick={() => setExpanded(true)} fluid>View {props.games.length - 5} more games...</Button>
        )];
    }
  }, [expanded, props.games]);

  if (props.games.length === 0) {
    return <p>No games present</p>;
  } else {
    return (
      <Item.Group divided>
        {games}
      </Item.Group>
    );
  }
}

function formatGame(game: Game, notes?: string): React.ReactNode {
  if (game.missing) {
    return (
      <Item>
        <Item.Image src='https://react.semantic-ui.com/images/wireframe/image.png'/>
        <Item.Content>
          <Item.Header>Missing Game</Item.Header>
          <Item.Description>Game could not be found</Item.Description>
        </Item.Content>
      </Item>
    );
  } else {
    const metaDesc = [game.developer, game.publisher].filter((s) => s !== '').join(' | ');
    return (
      <Item>
        <Item.Image src={`https://infinity.unstable.life/images/Logos/${game.id.slice(0,2)}/${game.id.slice(2,4)}/${game.id}.png`}/>
        <Item.Content>
          <Item.Header>{game.title}</Item.Header>
          <Item.Meta>{metaDesc}</Item.Meta>
          <Item.Description className='newline-preserved'>{notes ? `Playlist Note - ${notes}\n${game.originalDescription}` : game.originalDescription}</Item.Description>
          {game.filterGroups.length > 0 && (
            <Item.Extra>
              {game.filterGroups.map((group) => (
                <Label icon='warning sign' content={group}/>
              ))}
            </Item.Extra>
          )}
        </Item.Content>
      </Item>
    );
  }
}
