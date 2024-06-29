import ToggleTheme from '@/components/ToggleTheme';

const Navbar = () => {
  return (
    <div className="fixed top-0 w-full border-b border-b-accent">
      <div className="container flex items-center justify-between bg-primary-foreground px-2 py-2">
        <h1 className="text-2xl font-bold">Task Manager</h1>
        <ToggleTheme />
      </div>
    </div>
  );
};

export default Navbar;
