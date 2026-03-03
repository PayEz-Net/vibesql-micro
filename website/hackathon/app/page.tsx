import Navigation from '@/components/Navigation';
import Hero from '@/components/Hero';
import About from '@/components/About';
import Tracks from '@/components/Tracks';
import Prizes from '@/components/Prizes';
import Schedule from '@/components/Schedule';
import Rules from '@/components/Rules';
import Judges from '@/components/Judges';
import FAQ from '@/components/FAQ';
import Footer from '@/components/Footer';

export default function Home() {
  return (
    <main>
      <Navigation />
      <Hero />
      <About />
      <Tracks />
      <Prizes />
      <Schedule />
      <Rules />
      <Judges />
      <FAQ />
      <Footer />
    </main>
  );
}
