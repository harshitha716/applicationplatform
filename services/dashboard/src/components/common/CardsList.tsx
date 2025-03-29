import React from 'react';
import { cn } from 'utils/common';

type CardsListProps = {
  cards: any[]; // Array of data for each card
  CardComponent: React.FC<any>; // Custom component to render each card
  gridClassName?: string; // CSS Grid template for columns
};

const CardsList: React.FC<CardsListProps> = ({
  cards,
  CardComponent,
  gridClassName = 'grid-cols-[repeat(auto-fill,_minmax(100px,_1fr))] gap-4',
}) => {
  return (
    <div className={cn('grid', gridClassName)}>
      {cards?.map((cardData, index) => <CardComponent key={index} data={cardData} />)}
    </div>
  );
};

export default CardsList;
