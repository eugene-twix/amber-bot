import { Routes, Route } from 'react-router-dom';
import { TabBar } from '@/components/TabBar';
import { RatingPage } from '@/pages/RatingPage';
import { TeamsPage } from '@/pages/TeamsPage';
import { TeamDetailPage } from '@/pages/TeamDetailPage';
import { TournamentsPage } from '@/pages/TournamentsPage';
import { TournamentDetailPage } from '@/pages/TournamentDetailPage';
import { ManagePage } from '@/pages/ManagePage';

function App() {
  return (
    <div className="min-h-screen bg-background">
      <div className="max-w-2xl mx-auto px-4">
        <Routes>
          <Route path="/" element={<RatingPage />} />
          <Route path="/teams" element={<TeamsPage />} />
          <Route path="/teams/:id" element={<TeamDetailPage />} />
          <Route path="/tournaments" element={<TournamentsPage />} />
          <Route path="/tournaments/:id" element={<TournamentDetailPage />} />
          <Route path="/manage" element={<ManagePage />} />
        </Routes>
      </div>
      <TabBar />
    </div>
  );
}

export default App;
