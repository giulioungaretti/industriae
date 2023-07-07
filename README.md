# Coding Assignment  

Industrial Control System prototype coding assignment for the **Head of Software Development** [role](https://candidate.hr-manager.net/ApplicationInit.aspx?cid=2723&ProjectId=143571&DepartmentId=18965&MediaId=5&SkipAdvertisement=False) at [ATLANT 3D](https://www.atlant3d.com/).

## Objective:

Your task is to design and implement a prototype of a software system for a hypothetical Industrial Control System (ICS) of your choice. You may choose any of the following examples: an automated greenhouse system, an automated factory line, a smart traffic control system, or an automated water treatment plant. Alternatively, you may propose your own ICS scenario. 

The system should include a user interface for system control, an embedded system for real-time control of a single piece of industrial equipment (e.g., a conveyor belt in the case of an automated factory line), data logging for system monitoring, and safety protocols.

## Instructions

1. You may choose any programming language or technology stack for this assignment. In your submission, please provide a brief justification for your choices, explaining why you believe they are well-suited to the task at hand.

2. Document your architecture in a manner that is clear and easy to understand. This could involve writing detailed explanations of your design decisions, creating diagrams to visualize the structure of your system, using flowcharts to illustrate the flow of data and control in your system, or employing any other visual aids that you find helpful. The goal is to convey your ideas effectively to the reviewers, so feel free to use any approach that you think will achieve this.

## Requirements:

### User Interface 

Develop a simple UI that allows users to start and stop the equipment, and adjust a set of parameters (e.g., speed of the conveyor belt in the case of an automated factory line). The UI should also display real-time data from the equipment.

### Embedded System

Implement a simulated embedded system that controls the equipment based on the parameters set in the UI. The system should generate real-time data (e.g., current speed, status of the equipment in the case of an automated factory line).

### Data Logging

Implement a logging system that records the real-time data from the equipment and the parameters set by the user. The logs should be easily accessible for system monitoring.

### Safety Protocols

Implement basic safety protocols. For example, the system should automatically stop the equipment if the speed exceeds a certain limit.

### API Design

Design the APIs that would be used for communication between the UI, the embedded system, and the logging system.

### Error Handling 

The system should be able to handle potential errors or exceptional conditions. This could include sensor failures, communication errors, or any other issues that might arise.

### Testing Strategy 

Implement unit tests for your APIs and integration tests for your system as a whole.

## Evaluation Criteria

Your solution will be evaluated on the following criteria:

- **System Design**: Is your system design suitable for the problem at hand? Does it consider important factors such as latency, fault tolerance, and scalability?

- **Implementation**: Is your implementation consistent with your design? Does it meet the requirements of the assignment?

- **API Design**: Are your APIs well-designed and suitable for their intended purposes?

- **Error Handling**: Does your system have a robust strategy for handling errors and exceptional conditions?

- **Testing Strategy**: Is your testing strategy thorough and well-thought-out? Are your tests effective and comprehensive?

- **Code Quality**: Is your code clean, well-organized, and adhering to standard coding conventions?

## Submission Instructions

1. Fork this repository to your own GitHub account and make the forked repository private.

2. Clone the forked repository to your local machine.

3. Commit your changes to a new branch with the name `solution`.

4. Push your branch to the repository.

5. In your private repository, open a pull request from your branch to the `main` branch.

6. In the pull request description, provide a brief explanation of your solution and any decisions you made while developing it.

7. Finally, add the owner of this repository as a collaborator to your private repository so that we can review your solution.

**Please note**: Do not open a pull request in the original public repository. Your solution should be submitted in your private repository only.

Please submit your complete code files along with a README that explains how to run your prototype. Submit also your complete documentation along with any code that explains how to understand your architecture. Include any assumptions or design decisions you made while developing the architecture.
