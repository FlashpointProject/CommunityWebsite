import { UserPerm } from './utils/permissions';

export type User = {
  id: string;
  authed: boolean;
  username: string;
  avatarUrl: string;
  roles: string[];
  perms: UserPerm[];
};

export type Playlist = PlaylistInfo & {
  games: PlaylistGame[];
};

export type PlaylistInfo = {
  id: string;
  name: string;
  author: User;
  totalGames: number;
  library: string;
  description: string;
  extreme: boolean;
  filterGroups: string[];
};

export type Tag = unknown;

export type Game = {
  id: string;
  title: string;
  series: string;
  developer: string;
  publisher: string;
  releaseDate: string;
  playMode: string[];
  language: string[];
  originalDescription: string;
  platform: string;
  tags: Tag[];
  updatedAt: string;
  extreme: boolean;
  filterGroups: string[];
  missing: boolean;
};

export type PlaylistGame = {
  game: Game | null;
  gameId: string;
  notes: string;
};

export type RawTag = unknown;

export type RawGame = {
  id: string;
  title: string;
  series: string;
  developer: string;
  publisher: string;
  release_date: string;
  play_mode: string[];
  language: string[];
  original_description: string;
  platform: string;
  tags: RawTag[];
  extreme: boolean;
  filter_groups: string[];
  updated_at: string;
  missing: boolean;
};

export type RawPlaylistGame = {
  game: RawGame | null;
  game_id: string;
  notes: string;
};

export type RawPlaylist = RawPlaylistInfo & {
  games: RawPlaylistGame[];
};

export type RawPlaylistInfo = {
  id: string;
  name: string;
  author: RawUser;
  total_games: number;
  library: string;
  description: string;
  extreme: boolean;
  filter_groups: string[];
};

export type RawUser = {
  uid: string;
  username: string;
  avatar_url: string;
  roles: string[];
};

export type DiscordRole = {
  id: string;
  name: string;
  color: string;
};

export enum DiscordRoleIds {
  'Administrator' = '441043545735036929',
  'Moderator' = '442462642599231499',
  'Archivist' = '475413811394904074',
  'Developer' = '871819872408055828',
  'Mechanic' = '477773789724409861',
  'Hacker' = '442987546046103562',
  'Hunter' = '546894295962222622',
  'Tester' = '442988314480476170',
  'Curator' = '442665038642413569',
  'Editor' = '1128307753459392513',
  'VIP' = '453048638646648833',
  'Donator' = '1136941292031594587',
  '2021 Donator' = '921852025904435251',
  'International Manager' = '1126109274418987048',
  'Helper' = '914324576434028574',
  'Translator' = '703740452385325136',
  'Former Staff' = '666444964120494080',
  'Virus Destroyer' = '602334900623900672',
  'MOTAS Finder' = '636078914493480960',
  'Trial Curator' = '569328799318016018',
}

export type RawContentReport = {
  id: string;
  content_type: string;
  content_id: string;
  report_state: string;
  reported_user: RawUser;
  reported_by: RawUser;
  report_reason: string;
  resolved_at: string;
  resolved_by: RawUser;
  action_taken: string;
  created_at: string;
  updated_at: string;
};

export type ContentReport = {
  id: string;
  contentType: string;
  contentId: string;
  state: string;
  reportedUser: User;
  reportedBy: User;
  reportReason: string;
  resolvedAt: string;
  resolvedBy: User;
  actionTaken: string;
  createdAt: string;
  updatedAt: string;
};

export const RolesWithIcon: string[] = [
  DiscordRoleIds.Administrator
];

export type FilterPlaylists = {
  order: 'name' | 'total_games' | 'created_at' | 'updated_at';
  orderReverse: boolean;
  library: '' | 'arcade' | 'theatre';
  page: number;
  pageSize: number;
  extreme: boolean;
};

type ContentState = '' | 'reported' | 'resolved';
export type FilterContentReports = {
  order: 'created_at' | 'updated_at' | 'aggregate_report_score' | 'resolved_at' | 'action_taken' | 'content_type' | 'content_author' | 'reporter';
  content: '' | 'playlist' | 'comment';
  reportState: ContentState;
  orderReverse: boolean;
  page: number;
  pageSize: number;
};

export type NewsPost = {
  id: string;
  postType: string;
  title: string;
  content: string;
  author: User;
  createdAt: string;
  updatedAt: string;
};

export type RawNewsPost = {
  id: string;
  post_type: string;
  title: string;
  content: string;
  author: RawUser;
  created_at: string;
  updated_at: string;
};

export type GotdGame = {
  id: string;
  author: string;
  description: string;
  date: string;
};

export type GotdFile = {
  games: GotdGame[];
};

export type RawGotdSuggestion = {
  id: string;
  game: RawGame;
  author: string;
  description: string;
  suggested_date: string;
  created_at: string;
};

export type GotdSuggestion = {
  id: string;
  game: Game;
  author: string;
  description: string;
  suggestedDate: string;
  createdAt: string;
};

export type FilterGotdSuggestions = {
  order: 'suggested_date' | 'created_at';
  orderReverse: boolean;
  page: number;
  pageSize: number;
};
