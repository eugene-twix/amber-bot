import { useState } from 'react';
import { useMe, useTeams, useTournaments, useUsers, useCreateTeam, useCreateTournament, useCreateResult, useCreateMember, useUpdateUserRole } from '@/hooks/useApi';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Badge } from '@/components/ui/badge';
import { toast } from 'sonner';

export function ManagePage() {
  const { data: user } = useMe();
  const { data: teams } = useTeams();
  const { data: tournaments } = useTournaments();
  const { data: users } = useUsers();

  const createTeam = useCreateTeam();
  const createTournament = useCreateTournament();
  const createResult = useCreateResult();
  const createMember = useCreateMember();
  const updateUserRole = useUpdateUserRole();

  // Form states
  const [teamName, setTeamName] = useState('');
  const [tournamentName, setTournamentName] = useState('');
  const [tournamentDate, setTournamentDate] = useState('');
  const [tournamentLocation, setTournamentLocation] = useState('');
  const [selectedTeamForMember, setSelectedTeamForMember] = useState('');
  const [memberName, setMemberName] = useState('');
  const [selectedTournamentForResult, setSelectedTournamentForResult] = useState('');
  const [selectedTeamForResult, setSelectedTeamForResult] = useState('');
  const [resultPlace, setResultPlace] = useState('');

  const isOrgOrAdmin = user?.role === 'organizer' || user?.role === 'admin';
  const isAdmin = user?.role === 'admin';

  if (!isOrgOrAdmin) {
    return (
      <div className="p-4 text-center text-muted-foreground">
        Доступ запрещён
      </div>
    );
  }

  const handleCreateTeam = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!teamName.trim()) return;

    try {
      await createTeam.mutateAsync(teamName.trim());
      toast.success('Команда создана');
      setTeamName('');
    } catch {
      toast.error('Ошибка создания команды');
    }
  };

  const handleCreateTournament = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!tournamentName.trim() || !tournamentDate) return;

    try {
      await createTournament.mutateAsync({
        name: tournamentName.trim(),
        date: tournamentDate,
        location: tournamentLocation.trim(),
      });
      toast.success('Турнир создан');
      setTournamentName('');
      setTournamentDate('');
      setTournamentLocation('');
    } catch {
      toast.error('Ошибка создания турнира');
    }
  };

  const handleAddMember = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedTeamForMember || !memberName.trim()) return;

    try {
      await createMember.mutateAsync({
        teamId: Number(selectedTeamForMember),
        name: memberName.trim(),
      });
      toast.success('Участник добавлен');
      setMemberName('');
    } catch {
      toast.error('Ошибка добавления участника');
    }
  };

  const handleAddResult = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedTournamentForResult || !selectedTeamForResult || !resultPlace) return;

    try {
      await createResult.mutateAsync({
        tournamentId: Number(selectedTournamentForResult),
        teamId: Number(selectedTeamForResult),
        place: Number(resultPlace),
      });
      toast.success('Результат записан');
      setResultPlace('');
    } catch {
      toast.error('Ошибка записи результата');
    }
  };

  return (
    <div className="p-4 pb-20">
      <h1 className="text-xl font-semibold mb-4">Управление</h1>

      <Tabs defaultValue="team" className="w-full">
        <TabsList className="flex w-full mb-4 overflow-x-auto">
          <TabsTrigger value="team" className="flex-1 min-w-fit text-xs sm:text-sm px-2 sm:px-3">Команда</TabsTrigger>
          <TabsTrigger value="member" className="flex-1 min-w-fit text-xs sm:text-sm px-2 sm:px-3">Участник</TabsTrigger>
          <TabsTrigger value="tournament" className="flex-1 min-w-fit text-xs sm:text-sm px-2 sm:px-3">Турнир</TabsTrigger>
          <TabsTrigger value="result" className="flex-1 min-w-fit text-xs sm:text-sm px-2 sm:px-3">Результат</TabsTrigger>
          {isAdmin && <TabsTrigger value="users" className="flex-1 min-w-fit text-xs sm:text-sm px-2 sm:px-3">Права</TabsTrigger>}
        </TabsList>

        {/* Create Team */}
        <TabsContent value="team">
          <Card>
            <CardHeader className="p-4 pb-2">
              <CardTitle className="text-base">Создать команду</CardTitle>
            </CardHeader>
            <CardContent className="p-4 pt-2">
              <form onSubmit={handleCreateTeam} className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="teamName">Название</Label>
                  <Input
                    id="teamName"
                    value={teamName}
                    onChange={(e) => setTeamName(e.target.value)}
                    placeholder="Введите название команды"
                    maxLength={100}
                  />
                </div>
                <Button type="submit" disabled={createTeam.isPending || !teamName.trim()}>
                  {createTeam.isPending ? 'Создание...' : 'Создать'}
                </Button>
              </form>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Add Member */}
        <TabsContent value="member">
          <Card>
            <CardHeader className="p-4 pb-2">
              <CardTitle className="text-base">Добавить участника</CardTitle>
            </CardHeader>
            <CardContent className="p-4 pt-2">
              <form onSubmit={handleAddMember} className="space-y-4">
                <div className="space-y-2">
                  <Label>Команда</Label>
                  <Select value={selectedTeamForMember} onValueChange={setSelectedTeamForMember}>
                    <SelectTrigger>
                      <SelectValue placeholder="Выберите команду" />
                    </SelectTrigger>
                    <SelectContent>
                      {teams?.items.map((team) => (
                        <SelectItem key={team.id} value={String(team.id)}>
                          {team.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
                <div className="space-y-2">
                  <Label htmlFor="memberName">Имя участника</Label>
                  <Input
                    id="memberName"
                    value={memberName}
                    onChange={(e) => setMemberName(e.target.value)}
                    placeholder="Введите имя участника"
                    maxLength={100}
                  />
                </div>
                <Button
                  type="submit"
                  disabled={createMember.isPending || !selectedTeamForMember || !memberName.trim()}
                >
                  {createMember.isPending ? 'Добавление...' : 'Добавить'}
                </Button>
              </form>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Create Tournament */}
        <TabsContent value="tournament">
          <Card>
            <CardHeader className="p-4 pb-2">
              <CardTitle className="text-base">Создать турнир</CardTitle>
            </CardHeader>
            <CardContent className="p-4 pt-2">
              <form onSubmit={handleCreateTournament} className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="tournamentName">Название</Label>
                  <Input
                    id="tournamentName"
                    value={tournamentName}
                    onChange={(e) => setTournamentName(e.target.value)}
                    placeholder="Введите название турнира"
                    maxLength={100}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="tournamentDate">Дата</Label>
                  <Input
                    id="tournamentDate"
                    type="date"
                    value={tournamentDate}
                    onChange={(e) => setTournamentDate(e.target.value)}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="tournamentLocation">Место проведения</Label>
                  <Input
                    id="tournamentLocation"
                    value={tournamentLocation}
                    onChange={(e) => setTournamentLocation(e.target.value)}
                    placeholder="Введите место проведения"
                    maxLength={200}
                  />
                </div>
                <Button
                  type="submit"
                  disabled={createTournament.isPending || !tournamentName.trim() || !tournamentDate}
                >
                  {createTournament.isPending ? 'Создание...' : 'Создать'}
                </Button>
              </form>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Add Result */}
        <TabsContent value="result">
          <Card>
            <CardHeader className="p-4 pb-2">
              <CardTitle className="text-base">Записать результат</CardTitle>
            </CardHeader>
            <CardContent className="p-4 pt-2">
              <form onSubmit={handleAddResult} className="space-y-4">
                <div className="space-y-2">
                  <Label>Турнир</Label>
                  <Select value={selectedTournamentForResult} onValueChange={setSelectedTournamentForResult}>
                    <SelectTrigger>
                      <SelectValue placeholder="Выберите турнир" />
                    </SelectTrigger>
                    <SelectContent>
                      {tournaments?.items.map((tournament) => (
                        <SelectItem key={tournament.id} value={String(tournament.id)}>
                          {tournament.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
                <div className="space-y-2">
                  <Label>Команда</Label>
                  <Select value={selectedTeamForResult} onValueChange={setSelectedTeamForResult}>
                    <SelectTrigger>
                      <SelectValue placeholder="Выберите команду" />
                    </SelectTrigger>
                    <SelectContent>
                      {teams?.items.map((team) => (
                        <SelectItem key={team.id} value={String(team.id)}>
                          {team.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
                <div className="space-y-2">
                  <Label htmlFor="resultPlace">Место</Label>
                  <Input
                    id="resultPlace"
                    type="number"
                    min="1"
                    max="1000"
                    value={resultPlace}
                    onChange={(e) => setResultPlace(e.target.value)}
                    placeholder="Введите место"
                  />
                </div>
                <Button
                  type="submit"
                  disabled={
                    createResult.isPending ||
                    !selectedTournamentForResult ||
                    !selectedTeamForResult ||
                    !resultPlace
                  }
                >
                  {createResult.isPending ? 'Сохранение...' : 'Сохранить'}
                </Button>
              </form>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Admin: Manage User Roles */}
        {isAdmin && (
          <TabsContent value="users">
            <Card>
              <CardHeader className="p-4 pb-2">
                <CardTitle className="text-base">Управление правами</CardTitle>
              </CardHeader>
              <CardContent className="p-4 pt-2">
                <div className="space-y-3">
                  {users?.items.map((u) => (
                    <div
                      key={u.telegram_id}
                      className="flex items-center justify-between py-2 border-b last:border-0"
                    >
                      <div className="flex flex-col">
                        <span className="text-sm font-medium">
                          {u.username || `ID: ${u.telegram_id}`}
                        </span>
                        <Badge
                          variant={
                            u.role === 'admin'
                              ? 'default'
                              : u.role === 'organizer'
                              ? 'secondary'
                              : 'outline'
                          }
                          className="w-fit mt-1"
                        >
                          {u.role}
                        </Badge>
                      </div>
                      <Select
                        key={`${u.telegram_id}-${u.role}`}
                        value={u.role}
                        onValueChange={async (newRole) => {
                          try {
                            await updateUserRole.mutateAsync({
                              telegramId: u.telegram_id,
                              role: newRole,
                            });
                            toast.success('Роль обновлена');
                          } catch {
                            toast.error('Ошибка обновления роли');
                          }
                        }}
                      >
                        <SelectTrigger className="w-32">
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="viewer">viewer</SelectItem>
                          <SelectItem value="organizer">organizer</SelectItem>
                          <SelectItem value="admin">admin</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>
                  ))}
                  {users?.items.length === 0 && (
                    <p className="text-sm text-muted-foreground text-center py-4">
                      Нет пользователей
                    </p>
                  )}
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        )}
      </Tabs>
    </div>
  );
}
