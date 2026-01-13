import { useState, useMemo } from 'react';
import { useTournaments } from '@/hooks/useApi';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { Badge } from '@/components/ui/badge';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Link } from 'react-router-dom';

type SortOption = 'date_desc' | 'date_asc' | 'name_asc' | 'name_desc';

const sortOptions: { value: SortOption; label: string }[] = [
  { value: 'date_desc', label: 'Сначала новые' },
  { value: 'date_asc', label: 'Сначала старые' },
  { value: 'name_asc', label: 'По названию (А-Я)' },
  { value: 'name_desc', label: 'По названию (Я-А)' },
];

export function TournamentsPage() {
  const { data, isLoading, error } = useTournaments();
  const [sortBy, setSortBy] = useState<SortOption>('date_desc');

  const now = new Date();

  const sortedTournaments = useMemo(() => {
    if (!data?.items) return [];

    return [...data.items].sort((a, b) => {
      switch (sortBy) {
        case 'date_desc':
          return new Date(b.date).getTime() - new Date(a.date).getTime();
        case 'date_asc':
          return new Date(a.date).getTime() - new Date(b.date).getTime();
        case 'name_asc':
          return a.name.localeCompare(b.name, 'ru');
        case 'name_desc':
          return b.name.localeCompare(a.name, 'ru');
        default:
          return 0;
      }
    });
  }, [data?.items, sortBy]);

  if (error) {
    return (
      <div className="p-4 text-center text-red-500">
        Ошибка загрузки турниров
      </div>
    );
  }

  return (
    <div className="p-4 pb-20">
      <div className="flex items-center justify-between mb-4">
        <h1 className="text-xl font-semibold">Турниры</h1>
        <Select value={sortBy} onValueChange={(v) => setSortBy(v as SortOption)}>
          <SelectTrigger className="w-[180px]">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            {sortOptions.map((opt) => (
              <SelectItem key={opt.value} value={opt.value}>
                {opt.label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      {isLoading ? (
        <div className="space-y-3">
          {[...Array(5)].map((_, i) => (
            <Skeleton key={i} className="h-24 w-full" />
          ))}
        </div>
      ) : (
        <div className="space-y-3">
          {sortedTournaments.map((tournament) => {
            const tournamentDate = new Date(tournament.date);
            const isUpcoming = tournamentDate > now;

            return (
              <Link key={tournament.id} to={`/tournaments/${tournament.id}`}>
                <Card className="hover:bg-accent transition-colors">
                  <CardHeader className="p-4 pb-2">
                    <div className="flex items-center justify-between">
                      <CardTitle className="text-base">{tournament.name}</CardTitle>
                      {isUpcoming && (
                        <Badge variant="secondary">Предстоящий</Badge>
                      )}
                    </div>
                  </CardHeader>
                  <CardContent className="p-4 pt-0">
                    <div className="text-sm text-muted-foreground space-y-1">
                      <p>📅 {tournamentDate.toLocaleDateString('ru-RU')}</p>
                      {tournament.location && (
                        <p>📍 {tournament.location}</p>
                      )}
                    </div>
                  </CardContent>
                </Card>
              </Link>
            );
          })}
          {sortedTournaments.length === 0 && (
            <div className="text-center text-muted-foreground py-8">
              Нет турниров
            </div>
          )}
        </div>
      )}
    </div>
  );
}
