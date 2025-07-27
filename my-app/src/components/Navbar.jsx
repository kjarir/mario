import GooeyNav from './GooeyNav/GooeyNav'

// update with your own items
const items = [
  { label: "Home", href: "#" },
  { label: "About", href: "#" },
  { label: "Contact", href: "#" },
];

const NavbarComponent = () => {
  return (
    <div style={{ backgroundColor: '#D9272C', height: '60px', position: 'relative' }}>
        <GooeyNav
      items={items}
      particleCount={15}
      particleDistances={[90, 10]}
      particleR={100}
      initialActiveIndex={0}
      animationTime={600}
      timeVariance={300}
      colors={[1, 2, 3, 1, 2, 3, 1, 4]}
    />
    </div>
    
  );
};

export default NavbarComponent;