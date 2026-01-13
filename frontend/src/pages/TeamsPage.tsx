import { useState, useMemo } from 'react';
import { useTeams } from '@/hooks/useApi';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { Input } from '@/components/ui/input';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Link } from 'react-router-dom';

type SortOption = 'name_asc' | 'name_desc' | 'date_desc' | 'date_asc';

const sortOptions: { value: SortOption; label: string }[] = [
  { value: 'name_asc', label: 'По названию (А-Я)' },
  { value: 'name_desc', label: 'По названию (Я-А)' },
  { value: 'date_desc', label: 'Сначала новые' },
  { value: 'date_asc', label: 'Сначала старые' },
];

export function TeamsPage() {
  const { data, isLoading, error } = useTeams();
  const [sortBy, setSortBy] = useState<SortOption>('name_asc');
  const [search, setSearch] = useState('');

  const filteredAndSortedTeams = useMemo(() => {
    if (!data?.items) return [];

    const searchLower = search.toLowerCase().trim();
    const filtered = searchLower
      ? data.items.filter((team) => team.name.toLowerCase().includes(searchLower))
      : data.items;

    return [...filtered].sort((a, b) => {
      switch (sortBy) {
        case 'name_asc':
          return a.name.localeCompare(b.name, 'ru');
        case 'name_desc':
          return b.name.localeCompare(a.name, 'ru');
        case 'date_desc':
          return new Date(b.created_at).getTime() - new Date(a.created_at).getTime();
        case 'date_asc':
          return new Date(a.created_at).getTime() - new Date(b.created_at).getTime();
        default:
          return 0;
      }
    });
  }, [data?.items, sortBy, search]);

  if (error) {
    return (
      <div className="p-4 text-center text-red-500">
        Ошибка загрузки команд
      </div>
    );
  }

  return (
    <div className="p-4 pb-20">
      <div className="flex items-center justify-between mb-4">
        <h1 className="text-xl font-semibold">Команды</h1>
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

      <Input
        type="search"
        placeholder="Поиск по названию..."
        value={search}
        onChange={(e) => setSearch(e.target.value)}
        className="mb-4"
      />

      {isLoading ? (
        <div className="space-y-3">
          {[...Array(5)].map((_, i) => (
            <Skeleton key={i} className="h-20 w-full" />
          ))}
        </div>
      ) : (
        <div className="space-y-3">
          {filteredAndSortedTeams.map((team) => (
            <Link key={team.id} to={`/teams/${team.id}`}>
              <Card className="hover:bg-accent transition-colors">
                <CardHeader className="p-4 pb-2">
                  <CardTitle className="text-base">{team.name}</CardTitle>
                </CardHeader>
                <CardContent className="p-4 pt-0">
                  <p className="text-sm text-muted-foreground">
                    Создана: {new Date(team.created_at).toLocaleDateString('ru-RU')}
                  </p>
                </CardContent>
              </Card>
            </Link>
          ))}
          {filteredAndSortedTeams.length === 0 && (
            <div className="text-center text-muted-foreground py-8">
              Нет зарегистрированных команд
            </div>
          )}
        </div>
      )}
    </div>
  );
}
