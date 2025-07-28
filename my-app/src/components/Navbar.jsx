import GooeyNav from './GooeyNav/GooeyNav'

// update with your own items
const items = [
  { label: "Home", href: "#" },
  { label: "Services", href: "#" },
  { label: "About", href: "#" },
  { label: "Contact", href: "#" },
  { label: "Pricing", href: "#" },
];

const NavbarComponent = () => {
  return (
    <div style={{ backgroundColor: '#D9272C', height: '50px', position: 'relative', paddingTop: '15px' }}>
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