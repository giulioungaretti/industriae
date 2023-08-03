import React from "react";

interface LedProps {
  color: string;
  size: number;
  on: boolean;
}

const Led: React.FC<LedProps> = ({ color, size, on }) => {
  const style = {
    width: size,
    height: size,
    borderRadius: "50%",
    backgroundColor: on ? color : "red",
    margin: `${size / 20}em`,
    boxShadow: on
      ? `0 0 ${size / 10}em ${size / 20}em ${color},
		   0 0 ${size / 5}em ${size / 10}em rgba(255, 255, 255, 0.5)`
      : "none",
    border: `${size / 20}em solid rgba(255, 255, 255, 0.2)`,
  };

  return <div className="status-led" style={style} />;
};

export default Led;

// export default Led;
