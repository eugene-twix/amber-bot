import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '@/services/api';

// Query keys
export const queryKeys = {
  me: ['me'] as const,
  teams: ['teams'] as const,
  team: (id: number) => ['team', id] as const,
  teamMembers: (id: number) => ['team', id, 'members'] as const,
  teamResults: (id: number) => ['team', id, 'results'] as const,
  tournaments: ['tournaments'] as const,
  tournament: (id: number) => ['tournament', id] as const,
  tournamentResults: (id: number) => ['tournament', id, 'results'] as const,
  rating: ['rating'] as const,
  users: ['users'] as const,
};

// User
export function useMe() {
  return useQuery({
    queryKey: queryKeys.me,
    queryFn: api.getMe,
    retry: false,
  });
}

// Teams
export function useTeams() {
  return useQuery({
    queryKey: queryKeys.teams,
    queryFn: api.getTeams,
  });
}

export function useTeam(id: number) {
  return useQuery({
    queryKey: queryKeys.team(id),
    queryFn: () => api.getTeam(id),
    enabled: !!id,
  });
}

export function useTeamMembers(id: number) {
  return useQuery({
    queryKey: queryKeys.teamMembers(id),
    queryFn: () => api.getTeamMembers(id),
    enabled: !!id,
  });
}

export function useTeamResults(id: number) {
  return useQuery({
    queryKey: queryKeys.teamResults(id),
    queryFn: () => api.getTeamResults(id),
    enabled: !!id,
  });
}

// Tournaments
export function useTournaments() {
  return useQuery({
    queryKey: queryKeys.tournaments,
    queryFn: api.getTournaments,
  });
}

export function useTournament(id: number) {
  return useQuery({
    queryKey: queryKeys.tournament(id),
    queryFn: () => api.getTournament(id),
    enabled: !!id,
  });
}

export function useTournamentResults(id: number) {
  return useQuery({
    queryKey: queryKeys.tournamentResults(id),
    queryFn: () => api.getTournamentResults(id),
    enabled: !!id,
  });
}

// Rating
export function useRating() {
  return useQuery({
    queryKey: queryKeys.rating,
    queryFn: api.getRating,
  });
}

// Users (Admin)
export function useUsers() {
  return useQuery({
    queryKey: queryKeys.users,
    queryFn: api.getUsers,
  });
}

// Mutations
export function useCreateTeam() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (name: string) => api.createTeam(name),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.teams });
    },
  });
}

export function useCreateTournament() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: { name: string; date: string; location: string }) =>
      api.createTournament(data.name, data.date, data.location),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.tournaments });
    },
  });
}

export function useCreateResult() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: { tournamentId: number; teamId: number; place: number }) =>
      api.createResult(data.tournamentId, data.teamId, data.place),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.rating });
    },
  });
}

export function useCreateMember() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: { teamId: number; name: string }) =>
      api.createMember(data.teamId, data.name),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.teamMembers(variables.teamId) });
    },
  });
}

export function useUpdateUserRole() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: { telegramId: number; role: string }) =>
      api.updateUserRole(data.telegramId, data.role),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.users });
    },
  });
}

// Update/Delete Teams
export function useUpdateTeam() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: { id: number; name: string; version: number }) =>
      api.updateTeam(data.id, data.name, data.version),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.teams });
      queryClient.invalidateQueries({ queryKey: queryKeys.team(variables.id) });
    },
  });
}

export function useDeleteTeam() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: { id: number; version: number }) =>
      api.deleteTeam(data.id, data.version),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.teams });
      queryClient.invalidateQueries({ queryKey: queryKeys.rating });
    },
  });
}

// Update/Delete Members
export function useUpdateMember() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: { teamId: number; memberId: number; name: string; version: number }) =>
      api.updateMember(data.teamId, data.memberId, data.name, data.version),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.teamMembers(variables.teamId) });
    },
  });
}

export function useDeleteMember() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: { teamId: number; memberId: number; version: number }) =>
      api.deleteMember(data.teamId, data.memberId, data.version),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.teamMembers(variables.teamId) });
    },
  });
}

// Update/Delete Tournaments
export function useUpdateTournament() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: { id: number; name: string; date: string; location: string; version: number }) =>
      api.updateTournament(data.id, data.name, data.date, data.location, data.version),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.tournaments });
      queryClient.invalidateQueries({ queryKey: queryKeys.tournament(variables.id) });
    },
  });
}

export function useDeleteTournament() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: { id: number; version: number }) =>
      api.deleteTournament(data.id, data.version),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.tournaments });
    },
  });
}

// Update/Delete Results
export function useUpdateResult() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: { tournamentId: number; resultId: number; place: number; version: number }) =>
      api.updateResult(data.tournamentId, data.resultId, data.place, data.version),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.tournamentResults(variables.tournamentId) });
      queryClient.invalidateQueries({ queryKey: queryKeys.rating });
    },
  });
}

export function useDeleteResult() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: { tournamentId: number; resultId: number; version: number }) =>
      api.deleteResult(data.tournamentId, data.resultId, data.version),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.tournamentResults(variables.tournamentId) });
      queryClient.invalidateQueries({ queryKey: queryKeys.rating });
    },
  });
}
