import { useState } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import { useTeam, useTeamMembers, useTeamResults, useMe, useUpdateTeam, useDeleteTeam, useDeleteMember } from '@/hooks/useApi';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog';
import { toast } from 'sonner';

export function TeamDetailPage() {
  const { id } = useParams<{ id: string }>();
  const teamId = Number(id);
  const navigate = useNavigate();

  const { data: user } = useMe();
  const { data: team, isLoading: teamLoading, error: teamError } = useTeam(teamId);
  const { data: members, isLoading: membersLoading } = useTeamMembers(teamId);
  const { data: results, isLoading: resultsLoading } = useTeamResults(teamId);

  const updateTeam = useUpdateTeam();
  const deleteTeam = useDeleteTeam();
  const deleteMember = useDeleteMember();

  const [editDialogOpen, setEditDialogOpen] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [editName, setEditName] = useState('');
  const [memberToDelete, setMemberToDelete] = useState<{ id: number; name: string; version: number } | null>(null);

  const canManage = user?.role === 'organizer' || user?.role === 'admin';

  const handleOpenEdit = () => {
    setEditName(team?.name || '');
    setEditDialogOpen(true);
  };

  const handleSaveEdit = async () => {
    if (!team || !editName.trim()) return;
    try {
      await updateTeam.mutateAsync({ id: team.id, name: editName.trim(), version: team.version });
      toast.success('Команда обновлена');
      setEditDialogOpen(false);
    } catch {
      toast.error('Ошибка обновления');
    }
  };

  const handleDelete = async () => {
    if (!team) return;
    try {
      await deleteTeam.mutateAsync({ id: team.id, version: team.version });
      toast.success('Команда удалена');
      navigate('/teams');
    } catch {
      toast.error('Ошибка удаления');
    }
  };

  const handleDeleteMember = async () => {
    if (!memberToDelete) return;
    try {
      await deleteMember.mutateAsync({
        teamId,
        memberId: memberToDelete.id,
        version: memberToDelete.version,
      });
      toast.success('Участник удалён');
      setMemberToDelete(null);
    } catch {
      toast.error('Ошибка удаления участника');
    }
  };

  if (teamError) {
    return (
      <div className="p-4 text-center text-red-500">
        Ошибка загрузки команды
      </div>
    );
  }

  return (
    <div className="p-4 pb-20">
      <Link to="/teams">
        <Button variant="ghost" size="sm" className="mb-4">
          ← Назад
        </Button>
      </Link>

      {teamLoading ? (
        <Skeleton className="h-12 w-3/4 mb-4" />
      ) : (
        <div className="flex items-center justify-between mb-4">
          <h1 className="text-xl font-semibold">{team?.name}</h1>
          {canManage && (
            <div className="flex gap-2">
              <Button variant="outline" size="sm" onClick={handleOpenEdit}>
                ✏️
              </Button>
              <Button variant="destructive" size="sm" onClick={() => setDeleteDialogOpen(true)}>
                🗑️
              </Button>
            </div>
          )}
        </div>
      )}

      {/* Members Section */}
      <Card className="mb-4">
        <CardHeader className="p-4 pb-2">
          <CardTitle className="text-base">Участники</CardTitle>
        </CardHeader>
        <CardContent className="p-4 pt-0">
          {membersLoading ? (
            <div className="space-y-2">
              {[...Array(3)].map((_, i) => (
                <Skeleton key={i} className="h-6 w-full" />
              ))}
            </div>
          ) : members?.items.length === 0 ? (
            <p className="text-sm text-muted-foreground">Нет участников</p>
          ) : (
            <ul className="space-y-1">
              {members?.items.map((member) => (
                <li key={member.id} className="flex items-center justify-between text-sm py-1">
                  <span>{member.name}</span>
                  {canManage && (
                    <Button
                      variant="ghost"
                      size="sm"
                      className="h-6 w-6 p-0 text-muted-foreground hover:text-destructive"
                      onClick={() => setMemberToDelete({ id: member.id, name: member.name, version: member.version })}
                    >
                      ×
                    </Button>
                  )}
                </li>
              ))}
            </ul>
          )}
        </CardContent>
      </Card>

      {/* Results Section */}
      <Card>
        <CardHeader className="p-4 pb-2">
          <CardTitle className="text-base">Результаты</CardTitle>
        </CardHeader>
        <CardContent className="p-4 pt-0">
          {resultsLoading ? (
            <div className="space-y-2">
              {[...Array(3)].map((_, i) => (
                <Skeleton key={i} className="h-10 w-full" />
              ))}
            </div>
          ) : results?.items.length === 0 ? (
            <p className="text-sm text-muted-foreground">Нет результатов</p>
          ) : (
            <div className="space-y-2">
              {results?.items.map((result) => (
                <div
                  key={result.id}
                  className="flex justify-between items-center py-2 border-b last:border-0"
                >
                  <div className="flex flex-col">
                    <span className="text-sm font-medium">
                      {result.tournament_name}
                    </span>
                    <span className="text-xs text-muted-foreground">
                      {result.tournament_date
                        ? new Date(result.tournament_date).toLocaleDateString('ru-RU')
                        : ''}
                    </span>
                  </div>
                  <div className="text-lg font-semibold">
                    {result.place} место
                  </div>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Edit Dialog */}
      <Dialog open={editDialogOpen} onOpenChange={setEditDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Редактировать команду</DialogTitle>
          </DialogHeader>
          <div className="space-y-4 py-4">
            <div className="space-y-2">
              <Label htmlFor="editTeamName">Название</Label>
              <Input
                id="editTeamName"
                value={editName}
                onChange={(e) => setEditName(e.target.value)}
                maxLength={100}
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setEditDialogOpen(false)}>
              Отмена
            </Button>
            <Button onClick={handleSaveEdit} disabled={updateTeam.isPending || !editName.trim()}>
              {updateTeam.isPending ? 'Сохранение...' : 'Сохранить'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Delete Team Dialog */}
      <AlertDialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Удалить команду?</AlertDialogTitle>
            <AlertDialogDescription>
              Команда "{team?.name}" будет удалена. Это действие нельзя отменить.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Отмена</AlertDialogCancel>
            <AlertDialogAction onClick={handleDelete} className="bg-destructive text-destructive-foreground hover:bg-destructive/90">
              {deleteTeam.isPending ? 'Удаление...' : 'Удалить'}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* Delete Member Dialog */}
      <AlertDialog open={!!memberToDelete} onOpenChange={(open) => !open && setMemberToDelete(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Удалить участника?</AlertDialogTitle>
            <AlertDialogDescription>
              Участник "{memberToDelete?.name}" будет удалён из команды.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Отмена</AlertDialogCancel>
            <AlertDialogAction onClick={handleDeleteMember} className="bg-destructive text-destructive-foreground hover:bg-destructive/90">
              {deleteMember.isPending ? 'Удаление...' : 'Удалить'}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
