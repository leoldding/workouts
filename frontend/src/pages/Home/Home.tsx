import React from "react";
import "./Home.css";
import workouts from "../../data/workouts.json";

const Home: React.FC = () => {
    return (
        <>
            <div className="home-container">
                {workouts && workouts.map(({ id, name, date, notes, exercises }) => (
                    <div className="workout" key={id}>
                        <h1>{name} - {date}</h1>
                        <p>{notes}</p>
                        {exercises && exercises.map(({ exercise_id, exercise_name, exercise_notes, reps, values, unit }) => (
                           <div className="exercies" key={exercise_id}>
                           <h2>{exercise_name}</h2>
                           <p>{exercise_notes}</p>
                           <div>{reps}</div>
                           <div>{values} {unit}</div>
                           </div>
                        ))}
                    </div>
                ))}
            </div>
        </>
    )
};

export default Home;
