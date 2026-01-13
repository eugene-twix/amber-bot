import { useState } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import { useTournament, useTournamentResults, useMe, useUpdateTournament, useDeleteTournament, useDeleteResult } from '@/hooks/useApi';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
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

export function TournamentDetailPage() {
  const { id } = useParams<{ id: string }>();
  const tournamentId = Number(id);
  const navigate = useNavigate();

  const { data: user } = useMe();
  const { data: tournament, isLoading: tournamentLoading, error: tournamentError } = useTournament(tournamentId);
  const { data: results, isLoading: resultsLoading } = useTournamentResults(tournamentId);

  const updateTournament = useUpdateTournament();
  const deleteTournament = useDeleteTournament();
  const deleteResult = useDeleteResult();

  const [editDialogOpen, setEditDialogOpen] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [editName, setEditName] = useState('');
  const [editDate, setEditDate] = useState('');
  const [editLocation, setEditLocation] = useState('');
  const [resultToDelete, setResultToDelete] = useState<{ id: number; teamName: string; version: number } | null>(null);

  const canManage = user?.role === 'organizer' || user?.role === 'admin';

  const handleOpenEdit = () => {
    setEditName(tournament?.name || '');
    setEditDate(tournament?.date || '');
    setEditLocation(tournament?.location || '');
    setEditDialogOpen(true);
  };

  const handleSaveEdit = async () => {
    if (!tournament || !editName.trim() || !editDate) return;
    try {
      await updateTournament.mutateAsync({
        id: tournament.id,
        name: editName.trim(),
        date: editDate,
        location: editLocation.trim(),
        version: tournament.version,
      });
      toast.success('Турнир обновлён');
      setEditDialogOpen(false);
    } catch {
      toast.error('Ошибка обновления');
    }
  };

  const handleDelete = async () => {
    if (!tournament) return;
    try {
      await deleteTournament.mutateAsync({ id: tournament.id, version: tournament.version });
      toast.success('Турнир удалён');
      navigate('/tournaments');
    } catch {
      toast.error('Ошибка удаления');
    }
  };

  const handleDeleteResult = async () => {
    if (!resultToDelete) return;
    try {
      await deleteResult.mutateAsync({
        tournamentId,
        resultId: resultToDelete.id,
        version: resultToDelete.version,
      });
      toast.success('Результат удалён');
      setResultToDelete(null);
    } catch {
      toast.error('Ошибка удаления результата');
    }
  };

  if (tournamentError) {
    return (
      <div className="p-4 text-center text-red-500">
        Ошибка загрузки турнира
      </div>
    );
  }

  const now = new Date();
  const tournamentDate = tournament?.date ? new Date(tournament.date) : null;
  const isUpcoming = tournamentDate && tournamentDate > now;

  return (
    <div className="p-4 pb-20">
      <Link to="/tournaments">
        <Button variant="ghost" size="sm" className="mb-4">
          ← Назад
        </Button>
      </Link>

      {tournamentLoading ? (
        <Skeleton className="h-12 w-3/4 mb-4" />
      ) : (
        <div className="mb-4">
          <div className="flex items-center justify-between mb-2">
            <div className="flex items-center gap-2">
              <h1 className="text-xl font-semibold">{tournament?.name}</h1>
              {isUpcoming && (
                <Badge variant="secondary">Предстоящий</Badge>
              )}
            </div>
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
          <div className="text-sm text-muted-foreground space-y-1">
            {tournamentDate && (
              <p>📅 {tournamentDate.toLocaleDateString('ru-RU')}</p>
            )}
            {tournament?.location && (
              <p>📍 {tournament.location}</p>
            )}
          </div>
        </div>
      )}

      {/* Results Section */}
      <Card>
        <CardHeader className="p-4 pb-2">
          <CardTitle className="text-base">Результаты</CardTitle>
        </CardHeader>
        <CardContent className="p-4 pt-0">
          {resultsLoading ? (
            <div className="space-y-2">
              {[...Array(5)].map((_, i) => (
                <Skeleton key={i} className="h-10 w-full" />
              ))}
            </div>
          ) : results?.items.length === 0 ? (
            <p className="text-sm text-muted-foreground">
              {isUpcoming ? 'Турнир ещё не проведён' : 'Нет результатов'}
            </p>
          ) : (
            <div className="space-y-2">
              {results?.items
                .sort((a, b) => a.place - b.place)
                .map((result) => (
                  <div
                    key={result.id}
                    className="flex justify-between items-center py-2 border-b last:border-0"
                  >
                    <Link
                      to={`/teams/${result.team_id}`}
                      className="flex items-center gap-3 flex-1 hover:bg-accent transition-colors rounded px-2 -mx-2 py-1"
                    >
                      <span className={`text-lg font-bold w-8 ${
                        result.place === 1 ? 'text-yellow-500' :
                        result.place === 2 ? 'text-gray-400' :
                        result.place === 3 ? 'text-amber-600' : ''
                      }`}>
                        {result.place}
                      </span>
                      <span className="text-sm font-medium">
                        {result.team_name}
                      </span>
                    </Link>
                    {canManage && (
                      <Button
                        variant="ghost"
                        size="sm"
                        className="h-6 w-6 p-0 text-muted-foreground hover:text-destructive ml-2"
                        onClick={() => setResultToDelete({ id: result.id, teamName: result.team_name || '', version: result.version })}
                      >
                        ×
                      </Button>
                    )}
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
            <DialogTitle>Редактировать турнир</DialogTitle>
          </DialogHeader>
          <div className="space-y-4 py-4">
            <div className="space-y-2">
              <Label htmlFor="editTournamentName">Название</Label>
              <Input
                id="editTournamentName"
                value={editName}
                onChange={(e) => setEditName(e.target.value)}
                maxLength={200}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="editTournamentDate">Дата</Label>
              <Input
                id="editTournamentDate"
                type="date"
                value={editDate}
                onChange={(e) => setEditDate(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="editTournamentLocation">Место проведения</Label>
              <Input
                id="editTournamentLocation"
                value={editLocation}
                onChange={(e) => setEditLocation(e.target.value)}
                maxLength={200}
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setEditDialogOpen(false)}>
              Отмена
            </Button>
            <Button onClick={handleSaveEdit} disabled={updateTournament.isPending || !editName.trim() || !editDate}>
              {updateTournament.isPending ? 'Сохранение...' : 'Сохранить'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Delete Tournament Dialog */}
      <AlertDialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Удалить турнир?</AlertDialogTitle>
            <AlertDialogDescription>
              Турнир "{tournament?.name}" будет удалён. Это действие нельзя отменить.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Отмена</AlertDialogCancel>
            <AlertDialogAction onClick={handleDelete} className="bg-destructive text-destructive-foreground hover:bg-destructive/90">
              {deleteTournament.isPending ? 'Удаление...' : 'Удалить'}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* Delete Result Dialog */}
      <AlertDialog open={!!resultToDelete} onOpenChange={(open) => !open && setResultToDelete(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Удалить результат?</AlertDialogTitle>
            <AlertDialogDescription>
              Результат команды "{resultToDelete?.teamName}" будет удалён.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Отмена</AlertDialogCancel>
            <AlertDialogAction onClick={handleDeleteResult} className="bg-destructive text-destructive-foreground hover:bg-destructive/90">
              {deleteResult.isPending ? 'Удаление...' : 'Удалить'}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
