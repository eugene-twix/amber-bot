import { useState, useMemo } from 'react';
import { Link } from 'react-router-dom';
import { useRating } from '@/hooks/useApi';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { Badge } from '@/components/ui/badge';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { cn } from '@/lib/utils';

type SortKey = 'team_name' | 'top_places' | 'avg_place';
type SortDir = 'asc' | 'desc';

function SortIcon({ active, dir }: { active: boolean; dir: SortDir }) {
  return (
    <span className={cn('ml-1 inline-block', active ? 'opacity-100' : 'opacity-30')}>
      {dir === 'asc' ? '↑' : '↓'}
    </span>
  );
}

export function RatingPage() {
  const { data, isLoading, error } = useRating();
  const [sortKey, setSortKey] = useState<SortKey>('top_places');
  const [sortDir, setSortDir] = useState<SortDir>('desc');

  const sortedData = useMemo(() => {
    if (!data?.items) return [];

    return [...data.items].sort((a, b) => {
      let cmp = 0;
      switch (sortKey) {
        case 'team_name':
          cmp = a.team_name.localeCompare(b.team_name, 'ru');
          break;
        case 'top_places':
          cmp = a.top_places - b.top_places;
          break;
        case 'avg_place':
          cmp = a.avg_place - b.avg_place;
          break;
      }
      return sortDir === 'asc' ? cmp : -cmp;
    });
  }, [data?.items, sortKey, sortDir]);

  const handleSort = (key: SortKey) => {
    if (sortKey === key) {
      setSortDir(sortDir === 'asc' ? 'desc' : 'asc');
    } else {
      setSortKey(key);
      setSortDir(key === 'avg_place' ? 'asc' : 'desc');
    }
  };

  if (error) {
    return (
      <div className="p-4 text-center text-red-500">
        Ошибка загрузки рейтинга
      </div>
    );
  }

  return (
    <div className="p-4 pb-20 min-h-[calc(100vh-5rem)] flex flex-col">
      <Card className="flex-1 flex flex-col">
        <CardHeader className="pb-3">
          <CardTitle className="text-lg">Рейтинг команд</CardTitle>
        </CardHeader>
        <CardContent className="p-0 flex-1">
          {isLoading ? (
            <div className="space-y-2 p-4">
              {[...Array(5)].map((_, i) => (
                <Skeleton key={i} className="h-12 w-full" />
              ))}
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="w-10">#</TableHead>
                  <TableHead
                    className="cursor-pointer hover:bg-muted/50 select-none"
                    onClick={() => handleSort('team_name')}
                  >
                    Команда
                    <SortIcon active={sortKey === 'team_name'} dir={sortKey === 'team_name' ? sortDir : 'asc'} />
                  </TableHead>
                  <TableHead
                    className="w-16 text-center cursor-pointer hover:bg-muted/50 select-none"
                    onClick={() => handleSort('top_places')}
                  >
                    🏆
                    <SortIcon active={sortKey === 'top_places'} dir={sortKey === 'top_places' ? sortDir : 'desc'} />
                  </TableHead>
                  <TableHead
                    className="w-32 text-center cursor-pointer hover:bg-muted/50 select-none"
                    onClick={() => handleSort('avg_place')}
                  >
                    Среднее место
                    <SortIcon active={sortKey === 'avg_place'} dir={sortKey === 'avg_place' ? sortDir : 'asc'} />
                  </TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {sortedData.map((team, index) => (
                  <TableRow key={team.team_id}>
                    <TableCell className="font-medium">{index + 1}</TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        <Link
                          to={`/teams/${team.team_id}`}
                          className="hover:underline hover:text-primary"
                        >
                          {team.team_name}
                        </Link>
                        {index === 0 && sortKey === 'top_places' && sortDir === 'desc' && (
                          <Badge variant="default" className="text-xs">Лидер</Badge>
                        )}
                      </div>
                    </TableCell>
                    <TableCell className="text-center">
                      <span className="font-semibold">{team.top_places}</span>
                    </TableCell>
                    <TableCell className="text-center">
                      {team.avg_place.toFixed(2)}
                    </TableCell>
                  </TableRow>
                ))}
                {sortedData.length === 0 && (
                  <TableRow>
                    <TableCell colSpan={4} className="text-center text-muted-foreground py-8">
                      Нет данных о рейтинге
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
