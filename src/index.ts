/** Package imports */
import express from "express";
import * as bodyParser from "body-parser";

import "reflect-metadata";

import { globalRouter } from "./router/router";
import { Server } from "http";
import { Connection, createConnection } from "typeorm";

export const startServer = (port: string) => {
    return new Promise<{ server: Server; connection: Connection }>(
        (resolve, reject) => {
            /** Variables */
            const app: express.Application = express();

            /** Global middleware */
            app.use(bodyParser.json());

            /** Routes */
            app.use("/api", globalRouter);

            /** Start our server */
            createConnection({
                type: "postgres",
                host: process.env.DBHOST,
                port: Number(process.env.DBPORT),
                username: process.env.DBUSER,
                password: process.env.DBPASSWORD,
                database: process.env.DBDATABASE,
                synchronize: true,
                logging: false,
                entities: ["src/entity/**/*.ts"],
                migrations: ["src/migration/**/*.ts"],
                subscribers: ["src/subscriber/**/*.ts"],
            })
                .then((connection) => {
                    const server = app.listen(port, () => {
                        console.log(`Server is running on port ${port}...`);
                    });
                    resolve({ server, connection });
                })
                .catch((error: Error) => {
                    console.log(error);
                    reject();
                });
        }
    );
};

startServer(process.env.API_PORT);
