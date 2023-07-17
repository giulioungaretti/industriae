# industriae

A Industrial Control System prototype monorepo.

## Principles

Quality. All software must delight and work.
Safety.  All software must be safe.
Clarity. All software must have purpose and create value.

Based on this principles, a few choices were made:

*   No c/c++ unless documented exception or required (rtos/no llvm), if closer to metal is needed consider memory safe alternatives. E.g. rust.

*   The prototype will be scarce in features, but of high quality and clear purpose.

## Architecture

### The system
There is a pump controlled by an embedded rtos.
