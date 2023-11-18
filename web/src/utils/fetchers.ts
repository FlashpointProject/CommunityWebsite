import { Playlist, PlaylistGame, RawPlaylist } from '../types';
import { mapRawGame } from './mappers';
import { throwOn404 } from './permissions';

export async function fetchPlaylist(id: string): Promise<Playlist> {
  return fetch(`/api/playlist/${id}`)
  .then(throwOn404)
  .then((response) => response.json())
  .then((data) => {
    const rawPlaylist = data as RawPlaylist;
    // Map author to user
    const playlist: Playlist = {
      id: rawPlaylist.id,
      name: rawPlaylist.name,
      extreme: rawPlaylist.extreme,
      filterGroups: rawPlaylist.filter_groups,
      author: {
        authed: false,
        id: rawPlaylist.author.uid,
        username: rawPlaylist.author.username,
        avatarUrl: rawPlaylist.author.avatar_url,
        roles: rawPlaylist.author.roles,
        perms: [],
      },
      games: rawPlaylist.games.map<PlaylistGame>(rg => {
        if (rg.game === null) {
          return {
            game: null,
            gameId: rg.game_id,
            notes: rg.notes,
          };
        }
        return {
          game: mapRawGame(rg.game),
          gameId: rg.game_id,
          notes: rg.notes,
        };
      }),
      description: rawPlaylist.description,
      totalGames: rawPlaylist.total_games,
      library: rawPlaylist.library,
    };
    return playlist;
  });
}
