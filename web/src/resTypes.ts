import { PlaylistInfo, RawContentReport, RawGotdSuggestion, RawNewsPost } from './types';

export type ResponsePlaylists = {
  totalPlaylists: number;
  playlists: PlaylistInfo[];
};

export function mapResponsePlaylists(data: any): ResponsePlaylists {
  return {
    totalPlaylists: data.totalPlaylists,
    playlists: data.playlists
  };
}

export type ResponseUserProfile = {
  uid: string;
  username: string;
  avatar_url: string;
  roles: string[];
  updated_at: string;
};

export function mapResponseUserProfile(data: any): ResponseUserProfile {
  return {
    uid: data.uid,
    username: data.username,
    avatar_url: data.avatar_url,
    roles: data.roles,
    updated_at: data.updated_at
  };
}

export type ResponsePosts = {
  total: number;
  posts: RawNewsPost[];
};

export type ResponseContentReports = {
  total: number;
  reports: RawContentReport[];
};

export type ResponseGotdSuggestions = {
  total: number;
  suggestions: RawGotdSuggestion[];
};
