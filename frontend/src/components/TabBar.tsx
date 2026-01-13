import { Link, useLocation } from 'react-router-dom';
import { useMe } from '@/hooks/useApi';
import { cn } from '@/lib/utils';

const tabs = [
  { path: '/', label: 'Рейтинг', icon: '🏆' },
  { path: '/teams', label: 'Команды', icon: '📋' },
  { path: '/tournaments', label: 'Турниры', icon: '🎯' },
  { path: '/manage', label: 'Управление', icon: '⚙️', roles: ['organizer', 'admin'] },
];

export function TabBar() {
  const location = useLocation();
  const { data: user } = useMe();

  const visibleTabs = tabs.filter((tab) => {
    if (!tab.roles) return true;
    return user && tab.roles.includes(user.role);
  });

  return (
    <nav className="fixed bottom-0 left-0 right-0 bg-background border-t border-border">
      <div className="flex justify-around items-center h-16 max-w-2xl mx-auto px-4">
        {visibleTabs.map((tab) => {
          const isActive = tab.path === '/'
            ? location.pathname === '/'
            : location.pathname.startsWith(tab.path);

          return (
            <Link
              key={tab.path}
              to={tab.path}
              className={cn(
                'flex flex-col items-center justify-center flex-1 h-full gap-1',
                'text-muted-foreground hover:text-foreground transition-colors',
                isActive && 'text-primary'
              )}
            >
              <span className="text-xl">{tab.icon}</span>
              <span className="text-xs">{tab.label}</span>
            </Link>
          );
        })}
      </div>
    </nav>
  );
}
