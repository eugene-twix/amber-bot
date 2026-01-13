// API client for Amber Bot Mini App

const API_BASE = import.meta.env.VITE_API_URL || '/api/v1';

interface ListResponse<T> {
  items: T[];
  meta: {
    limit: number;
    offset: number;
    total: number;
  };
}

interface Team {
  id: number;
  name: string;
  created_at: string;
  version: number;
}

interface Member {
  id: number;
  name: string;
  team_id: number;
  joined_at: string;
  version: number;
}

interface Tournament {
  id: number;
  name: string;
  date: string;
  location: string;
  created_at: string;
  version: number;
}

interface Result {
  id: number;
  team_id: number;
  tournament_id: number;
  place: number;
  recorded_at: string;
  version: number;
  team_name?: string;
  tournament_name?: string;
  tournament_date?: string;
}

interface Rating {
  team_id: number;
  team_name: string;
  top_places: number;
  total_games: number;
  avg_place: number;
}

interface User {
  telegram_id: number;
  username: string;
  role: 'viewer' | 'organizer' | 'admin';
  created_at: string;
}

class ApiError extends Error {
  status: number;

  constructor(status: number, message: string) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
  }
}

// Get initData from Telegram Web App
function getInitData(): string {
  // @ts-expect-error Telegram WebApp global
  const tg = window.Telegram?.WebApp;
  return tg?.initData || '';
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const initData = getInitData();

  const response = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `TMA ${initData}`,
      ...options.headers,
    },
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'unknown_error' }));
    throw new ApiError(response.status, error.error || 'unknown_error');
  }

  return response.json();
}

// Public API
export const api = {
  // User
  getMe: () => request<User>('/public/me'),

  // Teams
  getTeams: () => request<ListResponse<Team>>('/public/teams'),
  getTeam: (id: number) => request<Team>(`/public/teams/${id}`),
  getTeamMembers: (id: number) => request<ListResponse<Member>>(`/public/teams/${id}/members`),
  getTeamResults: (id: number) => request<ListResponse<Result>>(`/public/teams/${id}/results`),

  // Tournaments
  getTournaments: () => request<ListResponse<Tournament>>('/public/tournaments'),
  getTournament: (id: number) => request<Tournament>(`/public/tournaments/${id}`),
  getTournamentResults: (id: number) => request<ListResponse<Result>>(`/public/tournaments/${id}/results`),

  // Rating
  getRating: () => request<ListResponse<Rating>>('/public/rating'),

  // Private - Teams
  createTeam: (name: string) => request<Team>('/private/teams', {
    method: 'POST',
    body: JSON.stringify({ name }),
  }),
  updateTeam: (id: number, name: string, version: number) => request<Team>(`/private/teams/${id}`, {
    method: 'PATCH',
    body: JSON.stringify({ name, version }),
  }),
  deleteTeam: (id: number, version: number) => request<{ deleted: boolean }>(`/private/teams/${id}`, {
    method: 'DELETE',
    body: JSON.stringify({ version }),
  }),

  // Private - Members
  createMember: (teamId: number, name: string) => request<Member>(`/private/teams/${teamId}/members`, {
    method: 'POST',
    body: JSON.stringify({ name }),
  }),
  updateMember: (teamId: number, memberId: number, name: string, version: number) => request<Member>(`/private/teams/${teamId}/members/${memberId}`, {
    method: 'PATCH',
    body: JSON.stringify({ name, version }),
  }),
  deleteMember: (teamId: number, memberId: number, version: number) => request<{ deleted: boolean }>(`/private/teams/${teamId}/members/${memberId}`, {
    method: 'DELETE',
    body: JSON.stringify({ version }),
  }),

  // Private - Tournaments
  createTournament: (name: string, date: string, location: string) => request<Tournament>('/private/tournaments', {
    method: 'POST',
    body: JSON.stringify({ name, date, location }),
  }),
  updateTournament: (id: number, name: string, date: string, location: string, version: number) => request<Tournament>(`/private/tournaments/${id}`, {
    method: 'PATCH',
    body: JSON.stringify({ name, date, location, version }),
  }),
  deleteTournament: (id: number, version: number) => request<{ deleted: boolean }>(`/private/tournaments/${id}`, {
    method: 'DELETE',
    body: JSON.stringify({ version }),
  }),

  // Private - Results
  createResult: (tournamentId: number, teamId: number, place: number) => request<Result>(`/private/tournaments/${tournamentId}/results`, {
    method: 'POST',
    body: JSON.stringify({ team_id: teamId, place }),
  }),
  updateResult: (tournamentId: number, resultId: number, place: number, version: number) => request<Result>(`/private/tournaments/${tournamentId}/results/${resultId}`, {
    method: 'PATCH',
    body: JSON.stringify({ place, version }),
  }),
  deleteResult: (tournamentId: number, resultId: number, version: number) => request<{ deleted: boolean }>(`/private/tournaments/${tournamentId}/results/${resultId}`, {
    method: 'DELETE',
    body: JSON.stringify({ version }),
  }),

  // Admin - Users
  getUsers: () => request<ListResponse<User>>('/private/users'),
  updateUserRole: (telegramId: number, role: string) => request<User>(`/private/users/${telegramId}/role`, {
    method: 'PUT',
    body: JSON.stringify({ role }),
  }),
};

export type { Team, Member, Tournament, Result, Rating, User, ListResponse };
export { ApiError };
