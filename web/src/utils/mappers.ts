import { ContentReport, Game, GotdSuggestion, NewsPost, RawContentReport, RawGame, RawGotdSuggestion, RawNewsPost } from '../types';

export function mapRawGame(data: RawGame): Game {
  return {
    id: data.id,
    title: data.title,
    series: data.series,
    developer: data.developer,
    publisher: data.publisher,
    platform: data.platform,
    playMode: data.play_mode,
    language: data.language,
    releaseDate: data.release_date,
    originalDescription: data.original_description,
    tags: data.tags,
    updatedAt: data.updated_at,
    extreme: data.extreme,
    filterGroups: data.filter_groups,
    missing: data.missing,
  };
}

export function mapRawNewsPost(data: RawNewsPost): NewsPost {
  return {
    id: data.id,
    postType: data.post_type,
    title: data.title,
    content: data.content,
    author: {
      id: data.author.uid,
      authed: true,
      username: data.author.username,
      avatarUrl: data.author.avatar_url,
      roles: data.author.roles,
      perms: [],
    },
    createdAt: data.created_at,
    updatedAt: data.updated_at
  };
}

export function mapRawContentReport(data: RawContentReport): ContentReport {
  return {
    id: data.id,
    contentRef: data.content_ref,
    state: data.report_state,
    reportedUser: {
      id: data.reported_user.uid,
      authed: true,
      username: data.reported_user.username,
      avatarUrl: data.reported_user.avatar_url,
      roles: data.reported_user.roles,
      perms: [],
    },
    reportedBy: {
      id: data.reported_by.uid,
      authed: true,
      username: data.reported_by.username,
      avatarUrl: data.reported_by.avatar_url,
      roles: data.reported_by.roles,
      perms: [],
    },
    reportReason: data.report_reason,
    context: data.context,
    resolvedAt: data.resolved_at,
    resolvedBy: data.resolved_by ? {
      id: data.resolved_by.uid,
      authed: true,
      username: data.resolved_by.username,
      avatarUrl: data.resolved_by.avatar_url,
      roles: data.resolved_by.roles,
      perms: [],
    } : undefined,
    actionTaken: data.action_taken,
    createdAt: data.created_at,
    updatedAt: data.updated_at
  };
}

export function mapRawGotdSuggestion(data: RawGotdSuggestion): GotdSuggestion {
  return {
    id: data.id,
    game: mapRawGame(data.game),
    author: data.author,
    description: data.description,
    suggestedDate: data.suggested_date,
    createdAt: data.created_at,
  };
}
